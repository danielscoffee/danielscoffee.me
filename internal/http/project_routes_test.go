package httpapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectIndexFiltersByTag(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/projects?tag=go", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Side Project") {
		t.Fatalf("expected matching project in body, got %q", body)
	}
	if strings.Contains(body, "No projects found") {
		t.Fatalf("expected non-empty project results, got %q", body)
	}
}
