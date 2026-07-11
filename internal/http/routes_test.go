package httpapp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
	"github.com/rs/zerolog"
)

func TestCompressionSkipsNoContent(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := testBlogServer().compressionMiddleware(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %x", w.Body.Bytes())
	}
	if got := w.Header().Get("Content-Encoding"); got != "" {
		t.Fatalf("expected no content encoding, got %q", got)
	}
}

func TestServerLimits(t *testing.T) {
	server := New(
		8080,
		"https://example.com",
		content.NewStore(nil),
		content.NewProjectStore(nil),
		content.Page{},
		zerolog.New(io.Discard),
		nil,
	)

	if server.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("expected 5s read header timeout, got %s", server.ReadHeaderTimeout)
	}
	if server.MaxHeaderBytes != 1<<20 {
		t.Fatalf("expected 1 MiB max headers, got %d", server.MaxHeaderBytes)
	}
}

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Page not found") {
		t.Fatalf("expected custom 404 page, got %q", w.Body.String())
	}
}

func TestRenderServerErrorReturnsCustomPage(t *testing.T) {
	s := testBlogServer()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/boom", nil)

	s.renderServerError(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Something broke") {
		t.Fatalf("expected custom 500 page, got %q", w.Body.String())
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
