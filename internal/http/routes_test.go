package httpapp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", w.Code)
	}
}

func TestRequestLoggingMiddleware_LogsRequestSummary(t *testing.T) {
	buf := &bytes.Buffer{}
	s := testBlogServer()
	s.logger = zerolog.New(buf)

	h := s.RegisterRoutes()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	got := buf.String()
	for _, marker := range []string{`"method":"GET"`, `"path":"/health"`, `"status":200`, `"remote_ip":"127.0.0.1"`, `"http_request"`} {
		if !strings.Contains(got, marker) {
			t.Fatalf("expected log to contain %q, got %s", marker, got)
		}
	}
}
