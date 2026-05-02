package httpapp

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

func testBlogServer() *Server {
	posts := []content.Post{
		{
			Title:    "Hello World",
			Slug:     "hello-world",
			Date:     "2026-04-26",
			Summary:  "Ship the first post",
			Tags:     []string{"go", "personal"},
			BodyMD:   "# Hello",
			BodyHTML: template.HTML(`<h1>Hello</h1>`),
		},
		{
			Title:    "Now",
			Slug:     "now",
			Date:     "2026-04-20",
			Summary:  "What I'm doing now",
			Tags:     []string{"now"},
			BodyMD:   "# Now",
			BodyHTML: template.HTML(`<h1>Now</h1>`),
		},
	}

	return &Server{
		port:         8080,
		contentStore: content.NewStore(posts),
		siteURL:      "https://example.com",
	}
}

func TestBlogRoutes(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path       string
		statusCode int
		contains   string
	}{
		{path: "/", statusCode: http.StatusOK, contains: "Daniel's Site"},
		{path: "/blog", statusCode: http.StatusOK, contains: "Blog"},
		{path: "/post/hello-world", statusCode: http.StatusOK, contains: "<article"},
		{path: "/tag/go", statusCode: http.StatusOK, contains: "Tagged with"},
		{path: "/post/missing", statusCode: http.StatusNotFound, contains: "404 page not found"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != tc.statusCode {
			t.Fatalf("path %s expected status %d got %d", tc.path, tc.statusCode, w.Code)
		}

		if !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatalf("path %s expected body to contain %q; got %q", tc.path, tc.contains, w.Body.String())
		}
	}
}

func TestBaseTemplateIncludesThemeControls(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	assertContainsAll(t, w.Body.String(), []string{
		`data-theme="system"`,
		`id="theme-toggle"`,
		`/assets/js/theme-init.js`,
		`/assets/js/theme-toggle.js`,
		`theme-preference`,
	})
}

func TestThemeAssetsAreServed(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path     string
		contains string
	}{
		{path: "/assets/js/theme-init.js", contains: "theme-preference"},
		{path: "/assets/js/theme-toggle.js", contains: "Theme:"},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))

		if w.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", tc.path, w.Code)
		}
		if !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatalf("%s expected body to contain %q", tc.path, tc.contains)
		}
	}
}

func TestPagesExposeStyleHooks(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path    string
		markers []string
	}{
		{
			path:    "/",
			markers: []string{"page-title", "page-subtitle", "section-title", "post-list", "post-item", "post-link", "tag-chip"},
		},
		{
			path:    "/blog",
			markers: []string{"page-title", "post-list", "post-meta-row"},
		},
		{
			path:    "/post/hello-world",
			markers: []string{"post-prose", "post-header", "post-title", "post-date"},
		},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))

		if w.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", tc.path, w.Code)
		}
		assertContainsAll(t, w.Body.String(), tc.markers)
	}
}

func assertContainsAll(t *testing.T, body string, markers []string) {
	t.Helper()
	for _, marker := range markers {
		if !strings.Contains(body, marker) {
			t.Fatalf("expected body to contain %q", marker)
		}
	}
}

func TestHealthRoute(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("expected content type application/json, got %q", got)
	}

	if !strings.Contains(w.Body.String(), `"status":"up"`) {
		t.Fatalf("expected health body to contain status up, got %q", w.Body.String())
	}
}
