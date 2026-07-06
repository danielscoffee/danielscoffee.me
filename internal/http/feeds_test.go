package httpapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedEndpoints(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path         string
		statusCode   int
		contentType  string
		bodyContains string
	}{
		{path: "/rss.xml", statusCode: http.StatusOK, contentType: "application/rss+xml", bodyContains: "<rss"},
		{path: "/sitemap.xml", statusCode: http.StatusOK, contentType: "application/xml", bodyContains: "<urlset"},
		{path: "/robots.txt", statusCode: http.StatusOK, contentType: "text/plain", bodyContains: "Sitemap:"},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))

		if w.Code != tc.statusCode {
			t.Fatalf("%s expected status %d got %d", tc.path, tc.statusCode, w.Code)
		}

		if got := w.Header().Get("Content-Type"); !strings.Contains(got, tc.contentType) {
			t.Fatalf("%s expected content type %q got %q", tc.path, tc.contentType, got)
		}

		if !strings.Contains(w.Body.String(), tc.bodyContains) {
			t.Fatalf("%s expected body to contain %q got %q", tc.path, tc.bodyContains, w.Body.String())
		}
	}
}

func TestFeedsIncludeProjectsAndSubposts(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	rss := httptest.NewRecorder()
	h.ServeHTTP(rss, httptest.NewRequest(http.MethodGet, "/rss.xml", nil))
	for _, marker := range []string{"/projects/side-project", "/projects/side-project/rebuild", "Side Project", "Rebuild Notes"} {
		if !strings.Contains(rss.Body.String(), marker) {
			t.Fatalf("expected RSS to contain %q, got %s", marker, rss.Body.String())
		}
	}

	sitemap := httptest.NewRecorder()
	h.ServeHTTP(sitemap, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))
	for _, marker := range []string{"/projects", "/projects/side-project", "/projects/side-project/rebuild"} {
		if !strings.Contains(sitemap.Body.String(), marker) {
			t.Fatalf("expected sitemap to contain %q, got %s", marker, sitemap.Body.String())
		}
	}
}
