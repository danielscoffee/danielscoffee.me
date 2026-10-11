package httpapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectDevlogList(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/projects/side-project", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	for _, marker := range []string{
		`<section class="section-block">`, `<h2 class="section-title">Devlog`,
		`<ul class="entry-list devlog-list">`,
		`<a class="entry-link devlog-link" href="/projects/side-project/rebuild">Rebuild Notes</a>`,
		`<span class="entry-meta"><time datetime="2026-05-10">May 10, 2026</time> · 1 min</span>`,
	} {
		if !strings.Contains(body, marker) {
			t.Errorf("expected devlog list markup %q", marker)
		}
	}
	for _, gone := range []string{"project-subpost-card", "how I rebuilt it", "year-mark"} {
		if strings.Contains(body, gone) {
			t.Errorf("flat devlog list should not contain %q", gone)
		}
	}
}

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
