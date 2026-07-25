package httpapp

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

func testBlogServer() *Server {
	posts := []content.Post{
		{
			Published: content.Published{
				Title:   "Hello World",
				Slug:    "hello-world",
				Date:    "2026-04-26",
				Summary: "Ship the first post",
				Tags:    []string{"go", "personal"},
			},
			BodyMD:   "# Hello",
			BodyHTML: template.HTML(`<h2>Hello</h2>`),
		},
		{
			Published: content.Published{
				Title:   "Now",
				Slug:    "now",
				Date:    "2026-04-20",
				Summary: "What I'm doing now",
				Tags:    []string{"now"},
			},
			BodyMD:   "# Now",
			BodyHTML: template.HTML(`<h2>Now</h2>`),
		},
	}

	projects := []content.Project{
		{
			Published: content.Published{
				Title:   "Side Project",
				Slug:    "side-project",
				Date:    "2026-05-01",
				Summary: "Small app",
				Tags:    []string{"go", "web"},
			},
			BodyMD:   "overview",
			BodyHTML: template.HTML(`<p>overview</p>`),
			SubPosts: []content.ProjectSubPost{
				{
					Published: content.Published{
						Title:   "Rebuild Notes",
						Slug:    "rebuild",
						Date:    "2026-05-10",
						Summary: "how I rebuilt it",
					},
					ParentSlug: "side-project",
					BodyMD:     "rebuild body",
					BodyHTML:   template.HTML(`<p>rebuild body</p>`),
				},
			},
		},
	}

	about := content.Page{
		Title:    "About Me",
		Slug:     "about",
		Date:     "2026-05-01",
		Summary:  "About",
		BodyMD:   "about text",
		BodyHTML: template.HTML(`<p>about text</p>`),
	}

	return &Server{
		port:          8080,
		contentStore:  content.NewStore(posts),
		projectStore:  content.NewProjectStore(projects),
		aboutPage:     about,
		siteURL:       "https://example.com",
		logger:        zerolog.New(io.Discard),
		searchIndexer: NewSearchIndexer(content.BuildSearchDocs(posts, projects)),
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
		{path: "/", statusCode: http.StatusFound, contains: ""},
		{path: "/blog", statusCode: http.StatusOK, contains: "Blog"},
		{path: "/about", statusCode: http.StatusOK, contains: "About Me"},
		{path: "/projects", statusCode: http.StatusOK, contains: "Projects"},
		{path: "/projects/side-project", statusCode: http.StatusOK, contains: "Side Project"},
		{path: "/projects/side-project", statusCode: http.StatusOK, contains: "project-subposts-grid"},
		{path: "/projects/side-project/rebuild", statusCode: http.StatusOK, contains: "Rebuild Notes"},
		{path: "/projects/side-project/missing", statusCode: http.StatusNotFound, contains: "Page not found"},
		{path: "/projects/missing", statusCode: http.StatusNotFound, contains: "Page not found"},
		{path: "/project/side-project", statusCode: http.StatusMovedPermanently, contains: ""},
		{path: "/post/hello-world", statusCode: http.StatusOK, contains: "<article"},
		{path: "/tag/go", statusCode: http.StatusOK, contains: "Tagged with"},
		{path: "/post/missing", statusCode: http.StatusNotFound, contains: "Page not found"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != tc.statusCode {
			t.Fatalf("path %s expected status %d got %d", tc.path, tc.statusCode, w.Code)
		}

		if tc.contains != "" && !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatalf("path %s expected body to contain %q; got %q", tc.path, tc.contains, w.Body.String())
		}
		if tc.path == "/" {
			if got := w.Header().Get("Location"); got != "/blog" {
				t.Fatalf("expected redirect to /blog, got %q", got)
			}
		}
	}
}

func TestSecurityAndCacheHeaders(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog", nil))

	if got := w.Header().Get("Content-Security-Policy"); !strings.Contains(got, "default-src 'self'") {
		t.Fatalf("expected CSP default-src self, got %q", got)
	}
	if got := w.Header().Get("Strict-Transport-Security"); got != "max-age=31536000; includeSubDomains" {
		t.Fatalf("unexpected HSTS header %q", got)
	}
	if got := w.Header().Get("Permissions-Policy"); got == "" {
		t.Fatal("expected Permissions-Policy header")
	}

	asset := httptest.NewRecorder()
	h.ServeHTTP(asset, httptest.NewRequest(http.MethodGet, "/assets/js/search.js", nil))
	if got := asset.Header().Get("Cache-Control"); got != "public, no-cache" {
		t.Fatalf("expected static assets to revalidate, got %q", got)
	}
}

func TestGzipCompression(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/blog", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip content encoding, got %q", got)
	}
	if got := w.Header().Get("Vary"); !strings.Contains(got, "Accept-Encoding") {
		t.Fatalf("expected Vary Accept-Encoding, got %q", got)
	}
}

func TestSiteShell(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	assertContainsAll(t, body, []string{
		`class="skip-link" href="#main-content"`,
		`<main id="main-content"`,
		`class="site-masthead"`,
		`<a href="/" class="site-brand">danielscoffee</a>`,
		`aria-label="Primary navigation"`,
		`<details class="mobile-menu">`,
		`<summary class="mobile-menu-toggle" aria-controls="mobile-menu-links header-actions">Menu</summary>`,
		`<nav id="mobile-menu-links" class="mobile-menu-links" aria-label="Mobile navigation">`,
		`<div id="header-actions" class="header-actions">`,
		`aria-label="Theme: System"`,
		`aria-label="Open search"`,
		`aria-labelledby="search-title"`,
		`<button class="search-close" type="submit">Close</button>`,
		`<footer class="site-footer"`,
		`class="footer-action-link" href="/rss.xml"`,
	})
	if strings.Contains(body, `class="site-tagline"`) {
		t.Error("site header should not contain tagline")
	}
}

func TestBaseTemplateIncludesThemeControls(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	assertContainsAll(t, w.Body.String(), []string{
		`data-theme="system"`,
		`id="theme-toggle"`,
		`href="/assets/css/output.css?v=2"`,
		`src="/assets/js/theme-init.js?v=2"`,
		`src="/assets/js/theme-toggle.js?v=2"`,
		`src="/assets/js/search.js?v=2"`,
		`theme-preference`,
		`search-modal`,
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
		{path: "/assets/js/search.js", contains: "Ctrl+K"},
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

func TestSearchRoute(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=projects+side", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"type":"projects"`) {
		t.Fatalf("expected projects result, got %s", w.Body.String())
	}
	if strings.Contains(w.Body.String(), `"type":"blog"`) {
		t.Fatalf("expected filtered results only, got %s", w.Body.String())
	}
}

func TestEditorialIndexes(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path    string
		markers []string
	}{
		{
			path: "/blog",
			markers: []string{
				`class="page-intro"`, `<h1 class="page-title">Blog</h1>`,
				`class="editorial-list post-list"`, `class="editorial-card post-item"`,
				`<time class="editorial-date" datetime="2026-04-26">2026-04-26</time>`,
				`href="/post/hello-world"`, `href="/tag/go"`,
			},
		},
		{
			path: "/projects",
			markers: []string{
				`class="page-intro"`, `<h1 class="page-title">Projects</h1>`,
				`class="editorial-list project-list"`, `class="editorial-card project-item"`,
				`<time class="editorial-date" datetime="2026-05-01">2026-05-01</time>`,
				`href="/projects/side-project"`, `href="/projects?tag=go"`,
			},
		},
		{
			path: "/tag/go",
			markers: []string{
				`class="page-intro"`, `Tagged with <span class="tag-emphasis">go</span>`,
				`class="editorial-list post-list"`, `class="editorial-card post-item"`,
				`datetime="2026-04-26"`, `href="/post/hello-world"`, `href="/tag/personal"`,
			},
		},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", tc.path, w.Code)
		}
		body := w.Body.String()
		assertContainsAll(t, body, tc.markers)
		if tc.path == "/tag/go" {
			heading := `<h2 class="section-title">Posts</h2>`
			headingIndex := strings.Index(body, heading)
			cardIndex := strings.Index(body, `class="editorial-card post-item"`)
			if headingIndex == -1 || cardIndex == -1 || headingIndex > cardIndex {
				t.Fatalf("%s expected %q before post cards", tc.path, heading)
			}
		}
	}
}

func TestEditorialArticles(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path       string
		date       string
		rawContent string
		breadcrumb string
	}{
		{path: "/post/hello-world", date: "2026-04-26", rawContent: `<h2>Hello</h2>`},
		{path: "/projects/side-project", date: "2026-05-01", rawContent: `<p>overview</p>`},
		{path: "/projects/side-project/rebuild", date: "2026-05-10", rawContent: `<p>rebuild body</p>`, breadcrumb: `href="/projects/side-project"`},
		{path: "/about", rawContent: `<p>about text</p>`},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", w.Code)
			}
			body := w.Body.String()
			assertContainsAll(t, body, []string{`<article class="article-shell`, `class="article-header`, `class="article-body`, tc.rawContent})
			if tc.date != "" {
				assertContainsAll(t, body, []string{`<time class="editorial-date" datetime="` + tc.date + `">` + tc.date + `</time>`})
			}
			if tc.breadcrumb != "" && !strings.Contains(body, tc.breadcrumb) {
				t.Fatalf("expected breadcrumb %q", tc.breadcrumb)
			}
			bodyStart := strings.Index(body, `class="article-body`)
			rawStart := strings.Index(body, tc.rawContent)
			if bodyStart == -1 || rawStart < bodyStart {
				t.Fatalf("expected raw content inside article body")
			}
		})
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
			path:    "/blog",
			markers: []string{"page-title", "page-subtitle", "section-title", "post-list", "post-item", "post-link", "tag-chip"},
		},
		{
			path:    "/blog",
			markers: []string{"page-title", "post-list", "post-meta-row"},
		},
		{
			path:    "/post/hello-world",
			markers: []string{"article-shell", "article-header", "article-title", "article-body", "editorial-date"},
		},
		{
			path:    "/projects",
			markers: []string{"project-list", "project-link"},
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
