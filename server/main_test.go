package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setup() {
	token = "testtoken-testtoken-testtoken-12345"
	origins = map[string]bool{"https://app.example.com": true}
	hosts = map[string]bool{"127.0.0.1:8787": true, "localhost:8787": true}
}

func do(method, path, body string, hdr map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://127.0.0.1:8787"+path, strings.NewReader(body))
	for k, v := range hdr {
		r.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	handle(w, r)
	return w
}

var okHdr = map[string]string{"Content-Type": "application/json", "Authorization": "Bearer testtoken-testtoken-testtoken-12345"}

func TestSecurityMatrix(t *testing.T) {
	setup()
	cases := []struct {
		name   string
		method string
		path   string
		body   string
		hdr    map[string]string
		want   int
	}{
		{"no token", "POST", "/exec", `{"cmd":"id"}`, map[string]string{"Content-Type": "application/json"}, http.StatusUnauthorized},
		{"wrong token", "POST", "/exec", `{"cmd":"id"}`, map[string]string{"Content-Type": "application/json", "Authorization": "Bearer nope"}, http.StatusUnauthorized},
		{"text/plain (CSRF)", "POST", "/exec", `{"cmd":"id"}`, map[string]string{"Content-Type": "text/plain", "Authorization": okHdr["Authorization"]}, http.StatusUnsupportedMediaType},
		{"foreign origin", "POST", "/exec", `{"cmd":"id"}`, map[string]string{"Content-Type": "application/json", "Authorization": okHdr["Authorization"], "Origin": "https://evil.com"}, http.StatusForbidden},
		{"foreign preflight", "OPTIONS", "/exec", "", map[string]string{"Origin": "https://evil.com"}, http.StatusForbidden},
		{"allowed preflight", "OPTIONS", "/exec", "", map[string]string{"Origin": "https://app.example.com"}, http.StatusNoContent},
		{"empty cmd", "POST", "/exec", `{"cmd":""}`, okHdr, http.StatusBadRequest},
		{"GET exec", "GET", "/exec", "", okHdr, http.StatusMethodNotAllowed},
		{"health", "GET", "/health", "", nil, http.StatusOK},
		{"test page", "GET", "/", "", nil, http.StatusOK},
		{"unknown path", "GET", "/nope", "", nil, http.StatusNotFound},
	}
	for _, c := range cases {
		if got := do(c.method, c.path, c.body, c.hdr).Code; got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

func TestBadHostRejected(t *testing.T) {
	setup()
	r := httptest.NewRequest("GET", "http://evil.com:8787/health", nil)
	w := httptest.NewRecorder()
	handle(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("DNS-rebinding host got %d, want 403", w.Code)
	}
}

func TestPrivateNetworkHeader(t *testing.T) {
	setup()
	w := do("OPTIONS", "/exec", "", map[string]string{"Origin": "https://app.example.com"})
	if w.Header().Get("Access-Control-Allow-Private-Network") != "true" {
		t.Fatal("missing Access-Control-Allow-Private-Network")
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "https://app.example.com" {
		t.Fatal("origin not echoed")
	}
}

func TestExecRunsAndTimesOut(t *testing.T) {
	setup()
	w := do("POST", "/exec", `{"cmd":"echo hi; echo err >&2; exit 3"}`, okHdr)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"errno":3`) || !strings.Contains(w.Body.String(), `"stdout":"hi\n"`) {
		t.Fatalf("unexpected: %d %s", w.Code, w.Body.String())
	}
	w = do("POST", "/exec", `{"cmd":"sleep 20","timeout":1}`, okHdr)
	if !strings.Contains(w.Body.String(), `"errno":124`) {
		t.Fatalf("timeout not enforced: %s", w.Body.String())
	}
}

// spawnResult collects a decoded /spawn response.
type spawnResult struct {
	out, err string
	code     int
	lines    int
}

func readSpawn(t *testing.T, body string) spawnResult {
	t.Helper()
	var res spawnResult
	res.code = -999
	for _, line := range strings.Split(strings.TrimSpace(body), "\n") {
		var ev struct {
			T    string `json:"t"`
			B    string `json:"b"`
			Code *int   `json:"code"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Fatalf("bad ndjson line %q: %v", line, err)
		}
		res.lines++
		switch ev.T {
		case "out", "err":
			raw, err := base64.StdEncoding.DecodeString(ev.B)
			if err != nil {
				t.Fatalf("bad base64: %v", err)
			}
			if ev.T == "out" {
				res.out += string(raw)
			} else {
				res.err += string(raw)
			}
		case "exit":
			res.code = *ev.Code
		}
	}
	return res
}

func TestSpawnRequiresAuthAndJSON(t *testing.T) {
	setup()
	if got := do("POST", "/spawn", `{"cmd":"id"}`, map[string]string{"Content-Type": "application/json"}).Code; got != http.StatusUnauthorized {
		t.Errorf("no token: got %d", got)
	}
	if got := do("POST", "/spawn", `{"cmd":"id"}`, map[string]string{"Content-Type": "text/plain", "Authorization": okHdr["Authorization"]}).Code; got != http.StatusUnsupportedMediaType {
		t.Errorf("text/plain: got %d", got)
	}
	foreign := map[string]string{"Content-Type": "application/json", "Authorization": okHdr["Authorization"], "Origin": "https://evil.com"}
	if got := do("POST", "/spawn", `{"cmd":"id"}`, foreign).Code; got != http.StatusForbidden {
		t.Errorf("foreign origin: got %d", got)
	}
	if got := do("GET", "/spawn", "", okHdr).Code; got != http.StatusMethodNotAllowed {
		t.Errorf("GET: got %d", got)
	}
}

func TestSpawnStreamsOutputAndExitCode(t *testing.T) {
	setup()
	w := do("POST", "/spawn", `{"cmd":"echo hi; echo oops >&2; exit 3"}`, okHdr)
	res := readSpawn(t, w.Body.String())
	if res.out != "hi\n" || res.err != "oops\n" || res.code != 3 {
		t.Fatalf("got %+v", res)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/x-ndjson" {
		t.Fatalf("content-type %q", ct)
	}
}

func TestSpawnArgsRunWithoutShell(t *testing.T) {
	setup()
	// With args the command is executed directly, so shell syntax stays literal.
	w := do("POST", "/spawn", `{"cmd":"echo","args":["a; echo injected","$HOME"]}`, okHdr)
	res := readSpawn(t, w.Body.String())
	if res.out != "a; echo injected $HOME\n" || res.code != 0 {
		t.Fatalf("got %+v", res)
	}
}

func TestSpawnMissingBinaryReportsError(t *testing.T) {
	setup()
	res := readSpawn(t, do("POST", "/spawn", `{"cmd":"/no/such/binary","args":[]}`, okHdr).Body.String())
	if res.code == 0 || res.code == -999 {
		t.Fatalf("expected failure exit code, got %+v", res)
	}
}

func TestEnvIsPassedAndValidated(t *testing.T) {
	setup()
	w := do("POST", "/exec", `{"cmd":"printf %s \"$FOO\"","env":{"FOO":"a b'c"}}`, okHdr)
	if !strings.Contains(w.Body.String(), `"stdout":"a b'c"`) {
		t.Fatalf("exec env: %s", w.Body.String())
	}
	res := readSpawn(t, do("POST", "/spawn", `{"cmd":"printf %s \"$FOO\"","env":{"FOO":"x"}}`, okHdr).Body.String())
	if res.out != "x" {
		t.Fatalf("spawn env: %+v", res)
	}
	if got := do("POST", "/exec", `{"cmd":"id","env":{"BAD NAME":"x"}}`, okHdr).Code; got != http.StatusBadRequest {
		t.Fatalf("invalid env name accepted: %d", got)
	}
}

func TestSpawnTimeoutKillsProcessGroup(t *testing.T) {
	setup()
	start := time.Now()
	res := readSpawn(t, do("POST", "/spawn", `{"cmd":"echo start; sleep 30 & sleep 30","timeout":1}`, okHdr).Body.String())
	if res.code != 124 || !strings.Contains(res.err, "timeout") || res.out != "start\n" {
		t.Fatalf("got %+v", res)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatalf("took too long: %v", time.Since(start))
	}
}

func TestSpawnKeepsSplitUTF8Intact(t *testing.T) {
	setup()
	// printf emits the 3 bytes of "€" in separate writes; base64 chunks must rejoin exactly.
	res := readSpawn(t, do("POST", "/spawn", `{"cmd":"printf '\\342'; sleep 0.2; printf '\\202\\254'"}`, okHdr).Body.String())
	if res.out != "€" {
		t.Fatalf("got %q", res.out)
	}
}
