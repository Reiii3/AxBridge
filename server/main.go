// Based on RiProG-ADB-Panel (MIT, (c) 2025 RiProG-id).
// Hardened: loopback bind, bearer token, Origin allowlist, CORS + Private
// Network Access preflight, timeout, size limits.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	maxBody    = 64 << 10
	maxOutput  = 1 << 20
	defTimeout = 30
	maxTimeout = 120

	// spawn streams, so it may run longer and produce more than exec.
	defSpawnTimeout = 300
	maxSpawnTimeout = 3600
	maxSpawnOutput  = 16 << 20
)

var envKey = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type execReq struct {
	Cmd     string            `json:"cmd"`
	Args    *[]string         `json:"args"` // spawn only: when present, run Cmd directly (no shell)
	Cwd     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
	Timeout int               `json:"timeout"`
}

type execRes struct {
	Errno  int    `json:"errno"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type capBuf struct {
	bytes.Buffer
	n int
}

func (c *capBuf) Write(p []byte) (int, error) {
	if c.n+len(p) > maxOutput {
		p = p[:maxOutput-c.n]
	}
	c.n += len(p)
	c.Buffer.Write(p)
	return len(p), nil // never fail the child on overflow
}

//go:embed index.html
var testPage []byte

var (
	token   string
	origins = map[string]bool{}
	hosts   = map[string]bool{}
)

func loadToken(path string) string {
	if t := os.Getenv("AXBRIDGE_TOKEN"); t != "" {
		return t
	}
	if b, err := os.ReadFile(path); err == nil && len(bytes.TrimSpace(b)) >= 32 {
		return string(bytes.TrimSpace(b))
	}
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		log.Fatal(err)
	}
	t := hex.EncodeToString(raw)
	if err := os.WriteFile(path, []byte(t+"\n"), 0o600); err != nil {
		log.Fatal(err)
	}
	return t
}

func cors(w http.ResponseWriter, r *http.Request) bool {
	o := r.Header.Get("Origin")
	if o == "" {
		return true // non-browser client (curl, local script)
	}
	if !origins[o] {
		return false
	}
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", o)
	h.Set("Vary", "Origin")
	h.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	h.Set("Access-Control-Allow-Private-Network", "true")
	h.Set("Access-Control-Max-Age", "600")
	return true
}

func authed(r *http.Request) bool {
	got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	return subtle.ConstantTimeCompare([]byte(got), []byte(token)) == 1
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Host check blocks DNS-rebinding: only our own loopback names are served.
	if !hosts[r.Host] {
		http.Error(w, "bad host", http.StatusForbidden)
		return
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if !cors(w, r) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	switch r.URL.Path {
	case "/", "/index.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(testPage)
	case "/health":
		writeJSON(w, 200, map[string]bool{"ok": true})
	case "/exec":
		doExec(w, r)
	case "/spawn":
		doSpawn(w, r)
	default:
		http.NotFound(w, r)
	}
}

// readRequest runs the checks shared by /exec and /spawn and decodes the body.
// It writes the error response itself and returns false when the request is refused.
func readRequest(w http.ResponseWriter, r *http.Request) (execReq, bool) {
	var req execReq
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return req, false
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return req, false
	}
	if !authed(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return req, false
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBody)).Decode(&req); err != nil || req.Cmd == "" {
		http.Error(w, "invalid json / empty cmd", http.StatusBadRequest)
		return req, false
	}
	for k := range req.Env {
		if !envKey.MatchString(k) {
			http.Error(w, "invalid env name", http.StatusBadRequest)
			return req, false
		}
	}
	return req, true
}

// buildCmd prepares the process: a shell line, or a direct exec when args were given.
func buildCmd(req execReq) *exec.Cmd {
	var cmd *exec.Cmd
	if req.Args != nil {
		cmd = exec.Command(req.Cmd, *req.Args...)
	} else {
		cmd = exec.Command("sh", "-c", req.Cmd)
	}
	cmd.Dir = req.Cwd
	if len(req.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range req.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

func clampTimeout(v, def, max int) int {
	if v <= 0 {
		return def
	}
	if v > max {
		return max
	}
	return v
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}

func doExec(w http.ResponseWriter, r *http.Request) {
	req, ok := readRequest(w, r)
	if !ok {
		return
	}
	http.NewResponseController(w).SetWriteDeadline(time.Now().Add((maxTimeout + 10) * time.Second))
	t := clampTimeout(req.Timeout, defTimeout, maxTimeout)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t)*time.Second)
	defer cancel()

	cmd := buildCmd(req)
	var so, se capBuf
	cmd.Stdout, cmd.Stderr = &so, &se
	if err := cmd.Start(); err != nil {
		writeJSON(w, 200, execRes{Errno: -1, Stderr: err.Error()})
		return
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	code := 0
	select {
	case err := <-done:
		code = exitCode(err)
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // whole process group
		<-done
		code = 124
		se.WriteString("\n[timeout after " + time.Duration(t*int(time.Second)).String() + "]")
	}
	writeJSON(w, 200, execRes{Errno: code, Stdout: so.String(), Stderr: se.String()})
}

// spawnEvent is one NDJSON line of the /spawn response.
//
//	{"t":"out","b":"<base64>"}   stdout chunk
//	{"t":"err","b":"<base64>"}   stderr chunk
//	{"t":"exit","code":N}        always the last line
//
// Chunks are base64 so a multi-byte character split across reads survives intact.
type spawnEvent struct {
	T    string `json:"t"`
	B    string `json:"b,omitempty"`
	Code *int   `json:"code,omitempty"`
}

func doSpawn(w http.ResponseWriter, r *http.Request) {
	req, ok := readRequest(w, r)
	if !ok {
		return
	}
	rc := http.NewResponseController(w)
	t := clampTimeout(req.Timeout, defSpawnTimeout, maxSpawnTimeout)
	rc.SetWriteDeadline(time.Now().Add(time.Duration(t+10) * time.Second))

	// r.Context() is cancelled when the client disconnects, which kills the process too.
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t)*time.Second)
	defer cancel()

	cmd := buildCmd(req)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		writeJSON(w, 200, spawnEvent{T: "exit", Code: ptr(-1)})
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	var (
		mu    sync.Mutex
		enc   = json.NewEncoder(w)
		total int
	)
	send := func(ev spawnEvent) {
		mu.Lock()
		defer mu.Unlock()
		enc.Encode(ev)
		rc.Flush()
	}
	pump := func(kind string, src io.Reader, wg *sync.WaitGroup) {
		defer wg.Done()
		buf := make([]byte, 32<<10)
		for {
			n, err := src.Read(buf)
			if n > 0 {
				mu.Lock()
				total += n
				over := total > maxSpawnOutput
				mu.Unlock()
				if over {
					syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
					return
				}
				send(spawnEvent{T: kind, B: base64.StdEncoding.EncodeToString(buf[:n])})
			}
			if err != nil {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go pump("out", stdout, &wg)
	go pump("err", stderr, &wg)

	// Watchdog: timeout or client disconnect kills the whole process group.
	finished := make(chan struct{})
	timedOut := false
	go func() {
		select {
		case <-ctx.Done():
			timedOut = ctx.Err() == context.DeadlineExceeded
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		case <-finished:
		}
	}()

	wg.Wait() // all output must be read before Wait closes the pipes
	err := cmd.Wait()
	close(finished)

	code := exitCode(err)
	if timedOut {
		code = 124
		send(spawnEvent{T: "err", B: base64.StdEncoding.EncodeToString([]byte("\n[timeout after " + (time.Duration(t) * time.Second).String() + "]"))})
	}
	send(spawnEvent{T: "exit", Code: &code})
}

func ptr[T any](v T) *T { return &v }

func main() {
	addr := flag.String("addr", "127.0.0.1:8787", "listen address (keep loopback)")
	allow := flag.String("origins", os.Getenv("AXBRIDGE_ORIGINS"),
		"comma-separated allowed browser origins, e.g. https://app.example.com")
	tokFile := flag.String("token-file", "token", "where the token is stored")
	pidFile := flag.String("pid-file", "", "write our pid here (optional)")
	flag.Parse()

	for _, o := range strings.Split(*allow, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins[strings.TrimRight(o, "/")] = true
		}
	}
	token = loadToken(*tokFile)
	if *pidFile != "" {
		os.WriteFile(*pidFile, []byte(strconv.Itoa(os.Getpid())+"\n"), 0o600)
	}

	_, port, _ := strings.Cut(*addr, ":")
	for _, h := range []string{"127.0.0.1:" + port, "localhost:" + port} {
		hosts[h] = true
		origins["http://"+h] = true // the built-in test page talks to itself
	}

	log.Printf("listening on http://%s  (token in %s)  origins=%d", *addr, *tokFile, len(origins))
	srv := &http.Server{
		Addr:              *addr,
		Handler:           http.HandlerFunc(handle),
		ReadHeaderTimeout: 5 * time.Second,
		// No global WriteTimeout: each handler sets its own deadline (spawn streams for longer).
	}
	log.Fatal(srv.ListenAndServe())
}
