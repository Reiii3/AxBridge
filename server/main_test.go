package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
