package httpapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectSubpostGrid(t *testing.T) {
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
		`<ul class="project-subposts-grid"`, `<li class="project-subpost-card">`,
		`<a class="project-subpost-link" href="/projects/side-project/rebuild">`,
		`<h3 class="project-subpost-title">Rebuild Notes</h3>`,
		`<p class="project-subpost-summary">how I rebuilt it</p>`,
		`<div class="project-subpost-date">`, `datetime="2026-05-10"`,
	} {
		if !strings.Contains(body, marker) {
			t.Errorf("expected project subpost markup %q", marker)
		}
	}
	if strings.Contains(body, "data-count") {
		t.Error("project subpost grid must not expose item count")
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
