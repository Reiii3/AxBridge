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
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	maxBody    = 64 << 10
	maxOutput  = 1 << 20
	defTimeout = 30
	maxTimeout = 120
)

type execReq struct {
	Cmd     string `json:"cmd"`
	Cwd     string `json:"cwd"`
	Timeout int    `json:"timeout"`
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
	default:
		http.NotFound(w, r)
	}
}

func doExec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	if !authed(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req execReq
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBody)).Decode(&req); err != nil || req.Cmd == "" {
		http.Error(w, "invalid json / empty cmd", http.StatusBadRequest)
		return
	}
	t := req.Timeout
	if t <= 0 {
		t = defTimeout
	}
	if t > maxTimeout {
		t = maxTimeout
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t)*time.Second)
	defer cancel()

	cmd := exec.Command("sh", "-c", req.Cmd)
	cmd.Dir = req.Cwd
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			code = -1
		}
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // whole process group
		<-done
		code = 124
		se.WriteString("\n[timeout after " + time.Duration(t*int(time.Second)).String() + "]")
	}
	writeJSON(w, 200, execRes{Errno: code, Stdout: so.String(), Stderr: se.String()})
}

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
		WriteTimeout:      (maxTimeout + 10) * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
