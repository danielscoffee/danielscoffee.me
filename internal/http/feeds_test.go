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
		{path: "/llms.txt", statusCode: http.StatusOK, contentType: "text/plain", bodyContains: "# danielscoffee"},
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

func TestRobotsTxt(t *testing.T) {
	want := "User-agent: *\n" +
		"Allow: /\n" +
		"Disallow: /search\n" +
		"Disallow: /health\n" +
		"\n" +
		"# Site map for language models and agents: https://example.com/llms.txt\n" +
		"Sitemap: https://example.com/sitemap.xml\n"
	if got := buildRobots("https://example.com/"); got != want {
		t.Fatalf("robots.txt mismatch:\n%s", got)
	}
}

func TestLLMSTxt(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/llms.txt", nil))
	body := w.Body.String()

	if !strings.HasPrefix(body, "# danielscoffee\n\n> ") {
		t.Fatalf("llms.txt must open with the H1 and a blockquote summary: %q", body[:min(len(body), 80)])
	}
	for _, want := range []string{
		"\n## Blog\n\n- [Hello World](https://example.com/post/hello-world): Ship the first post\n",
		"- [Now](https://example.com/post/now): What I'm doing now\n",
		"\n## Projects\n\n- [Side Project](https://example.com/projects/side-project): Small app\n",
		"- [Rebuild Notes](https://example.com/projects/side-project/rebuild): Side Project devlog: how I rebuilt it\n",
		"\n## About\n\n- [About Me](https://example.com/about): About\n",
		"\n## Optional\n\n- [RSS feed](https://example.com/rss.xml)",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("llms.txt missing %q in:\n%s", want, body)
		}
	}
}

func TestLLMSLinkEscapesBrackets(t *testing.T) {
	var b strings.Builder
	writeLLMSLink(&b, "Arrays [and] slices", "https://example.com/post/a", "note\nline")
	if got, want := b.String(), "- [Arrays \\[and\\] slices](https://example.com/post/a): note line\n"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
