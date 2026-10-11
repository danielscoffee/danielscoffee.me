package httpapp

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
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
		{path: "/", statusCode: http.StatusOK, contains: "danielscoffee"},
		{path: "/no-such-page", statusCode: http.StatusNotFound, contains: "Page not found"},
		{path: "/blog", statusCode: http.StatusOK, contains: "Blog"},
		{path: "/about", statusCode: http.StatusOK, contains: "About Me"},
		{path: "/projects", statusCode: http.StatusOK, contains: "Projects"},
		{path: "/projects/side-project", statusCode: http.StatusOK, contains: "Side Project"},
		{path: "/projects/side-project", statusCode: http.StatusOK, contains: "devlog-list"},
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
	if got := asset.Header().Get("Cache-Control"); got != "public, max-age=0, s-maxage=86400" {
		t.Fatalf("expected static assets cached at the edge and revalidated by browsers, got %q", got)
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
		`<a href="/blog" class="site-section" data-section="blog">blog</a>`,
		`<a href="/" class="site-back" aria-label="Back to home">`,
		`<div id="header-actions" class="header-actions">`,
		`href="https://github.com/danielscoffee"`,
		`href="https://www.linkedin.com/in/daniel-reis-santos/"`,
		`id="search-trigger"`,
		`aria-label="Theme"`,
		`aria-label="Open search"`,
		`aria-labelledby="search-title"`,
		`<button class="search-close" type="submit">Close</button>`,
		`<footer class="site-footer"`,
		`class="footer-action-link" href="/rss.xml"`,
	})
	if strings.Contains(body, `class="site-tagline"`) {
		t.Error("site header should not contain tagline")
	}
	for _, removed := range []string{"mobile-menu", "site-nav-desktop", "site-menu", `aria-label="Primary navigation"`, `class="social-link"`} {
		if strings.Contains(body, removed) {
			t.Errorf("site header should not contain %q", removed)
		}
	}
}

func TestHomePage(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if got := w.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("expected text/html, got %q", got)
	}

	body := w.Body.String()
	assertContainsAll(t, body, []string{
		`<h1 id="board-title" class="board-title">danielscoffee</h1>`,
		`<span class="door-name" data-section="blog">Blog</span>`,
		`<span class="door-name" data-section="projects">Projects</span>`,
		`<span class="door-name" data-section="about">About</span>`,
		`Backend developer. Systems enjoyer. Production bug survivor.`,
		`<a class="board-door" href="/blog">`,
		`<a class="board-door" href="/projects">`,
		`<a class="board-door" href="/about">`,
		`src="/assets/img/scenes/cafe-table.jpg"`,
		`src="/assets/img/scenes/pour-over.jpg"`,
		`src="/assets/img/scenes/takeaway-cup.jpg"`,
		`src="/assets/js/home-intro.js?v=1"`,
		`src="/assets/js/theme-init.js?v=2"`,
		`href="/assets/css/output.css?v=2"`,
	})
	blog := strings.Index(body, `href="/blog"`)
	projects := strings.Index(body, `href="/projects"`)
	about := strings.Index(body, `href="/about"`)
	if !(blog < projects && projects < about) {
		t.Errorf("expected doors in order Blog, Projects, About")
	}
	for _, absent := range []string{`class="site-masthead"`, `class="site-footer"`, `id="search-modal"`} {
		if strings.Contains(body, absent) {
			t.Errorf("home page should not contain %s", absent)
		}
	}
	for _, tag := range strings.Split(body, "<script")[1:] {
		end := strings.Index(tag, ">")
		if end == -1 || !strings.Contains(tag[:end], "src=") {
			t.Errorf("home page must not contain inline scripts: <script%s", tag)
		}
	}
}

func TestHomeAssetsAreServed(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/assets/js/home-intro.js", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "sessionStorage") {
		t.Fatalf("home-intro.js: status %d, body %q", w.Code, w.Body.String())
	}

	for _, scene := range []string{"cafe-table", "pour-over", "takeaway-cup"} {
		path := "/assets/img/scenes/" + scene + ".jpg"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", path, w.Code)
		}
		if got := w.Header().Get("Content-Type"); !strings.HasPrefix(got, "image/jpeg") {
			t.Errorf("%s expected image/jpeg, got %q", path, got)
		}
		body := w.Body.Bytes()
		if len(body) < 3 || body[0] != 0xFF || body[1] != 0xD8 || body[2] != 0xFF {
			t.Errorf("%s is not a JPEG", path)
		}
		if len(body) > 60*1024 {
			t.Errorf("%s is %d bytes; scene photos must stay under 60KB", path, len(body))
		}
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
		`src="/assets/js/code-copy.js?v=1"`,
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
				`<h1 class="sr-only">Blog</h1>`,
				`<div class="year-list post-list">`, `<h2 class="year-mark" id="year-2026">2026</h2>`,
				`<li class="entry post-item"><a class="entry-link post-link" href="/post/hello-world">Hello World</a>`,
				`<span class="entry-meta"><time datetime="2026-04-26">Apr 26</time> · 1 min</span>`,
			},
		},
		{
			path: "/projects",
			markers: []string{
				`<h1 class="sr-only">Projects</h1>`,
				`<div class="year-list project-list">`,
				`<a class="entry-link project-link" href="/projects/side-project">Side Project</a>`,
				`<span class="entry-meta"><time datetime="2026-05-01">May 1</time></span>`,
			},
		},
		{
			path: "/tag/go",
			markers: []string{
				`Tagged with <span class="tag-emphasis">go</span>`,
				`<div class="year-list post-list">`, `href="/post/hello-world"`,
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
		// Flat lists: no summaries, subtitles, tag chips, or cards.
		for _, gone := range []string{"Ship the first post", "Small app", `class="page-subtitle"`, "tag-chip", "editorial-card"} {
			if strings.Contains(body, gone) {
				t.Errorf("%s should not contain %q", tc.path, gone)
			}
		}
	}
}

func TestArticleTags(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	for path, want := range map[string]string{
		"/post/hello-world":      `<a class="article-tag" href="/tag/go">go</a>`,
		"/projects/side-project": `<a class="article-tag" href="/projects?tag=go">go</a>`,
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		assertContainsAll(t, w.Body.String(), []string{`<p class="article-meta">`, want})
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
			assertContainsAll(t, body, []string{`<article class="article-shell`, `class="article-body`, tc.rawContent})
			if tc.path == "/about" {
				// The masthead names the section, so About's heading is screen-reader only.
				assertContainsAll(t, body, []string{`<h1 class="sr-only">About Me</h1>`})
				if strings.Contains(body, `class="article-header`) {
					t.Fatal("about page should not repeat a visible title under the masthead")
				}
			} else {
				assertContainsAll(t, body, []string{`class="article-header`})
			}
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
			markers: []string{"year-list", "year-mark", "post-list", "post-item", "post-link", "entry-meta"},
		},
		{
			path:    "/blog",
			markers: []string{"entry-link", "entry-meta"},
		},
		{
			path:    "/post/hello-world",
			markers: []string{"article-shell", "article-header", "article-title", "article-body", "article-meta", "editorial-date"},
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

func TestPageMetadata(t *testing.T) {
	s := testBlogServer()
	s.canonicalBase = "https://example.com"
	h := s.RegisterRoutes()

	get := func(path string) string {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		return w.Body.String()
	}

	t.Run("share and canonical tags", func(t *testing.T) {
		assertContainsAll(t, get("/blog"), []string{
			`<title>Blog · danielscoffee</title>`,
			`<meta name="description" content="Every post on danielscoffee, newest first.">`,
			`<link rel="canonical" href="https://example.com/blog">`,
			`<meta property="og:url" content="https://example.com/blog">`,
			`<meta property="og:image" content="https://example.com/assets/img/og-default.jpg">`,
			`<meta property="og:title" content="Blog · danielscoffee">`,
			`<meta property="og:description" content="Every post on danielscoffee, newest first.">`,
			`<meta name="twitter:card" content="summary_large_image">`,
		})
	})

	t.Run("filtered projects name their tag and point home", func(t *testing.T) {
		body := get("/projects?tag=go")
		assertContainsAll(t, body, []string{
			`<title>Projects tagged go · danielscoffee</title>`,
			`content="Projects on danielscoffee tagged go, with their devlogs."`,
			`Tagged with <span class="tag-emphasis">go</span>`,
			`<link rel="canonical" href="https://example.com/projects">`,
		})
	})

	t.Run("tag pages have their own title and description", func(t *testing.T) {
		assertContainsAll(t, get("/tag/go"), []string{
			`<title>Posts tagged go · danielscoffee</title>`,
			`content="Posts on danielscoffee tagged go."`,
		})
	})

	t.Run("posts carry BlogPosting data", func(t *testing.T) {
		assertContainsAll(t, get("/post/hello-world"), []string{
			`<script id="article-data" type="application/ld+json">`,
			`"@type":"BlogPosting"`,
			`"headline":"Hello World"`,
			`"datePublished":"2026-04-26"`,
			`"url":"https://example.com/post/hello-world"`,
		})
	})

	t.Run("error pages are not shareable", func(t *testing.T) {
		body := get("/no-such-page")
		for _, tag := range []string{`rel="canonical"`, `property="og:`, `twitter:card`} {
			if strings.Contains(body, tag) {
				t.Errorf("404 page should not carry %s", tag)
			}
		}
	})
}

func TestPageMetadataWithoutSiteURL(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/blog", nil))
	body := w.Body.String()
	if strings.Contains(body, `rel="canonical"`) || strings.Contains(body, `property="og:url"`) || strings.Contains(body, `property="og:image"`) {
		t.Fatal("without SITE_URL the head must not guess absolute URLs")
	}
	assertContainsAll(t, body, []string{`<meta property="og:title" content="Blog · danielscoffee">`, `<meta name="twitter:card" content="summary_large_image">`})
}

func TestFavicon(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	for _, path := range []string{"/", "/blog"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		assertContainsAll(t, w.Body.String(), []string{`<link rel="icon" href="/assets/favicon.svg" type="image/svg+xml">`})
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/assets/favicon.svg", nil))
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "image/svg+xml") {
		t.Fatalf("favicon: status %d, type %q", w.Code, w.Header().Get("Content-Type"))
	}
	if body := w.Body.String(); !strings.Contains(body, `viewBox="0 0 32 32"`) || strings.Contains(body, "<script") {
		t.Fatalf("favicon must be a plain 32x32 SVG: %s", body)
	}
}

func TestBackLinks(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	for path, want := range map[string]string{
		"/projects/side-project":         `<a class="page-back" href="/projects" aria-label="Back">`,
		"/projects/side-project/rebuild": `<a class="page-back" href="/projects/side-project" aria-label="Back">`,
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		body := w.Body.String()
		assertContainsAll(t, body, []string{`<div class="page-breadcrumb">`, want, `<script defer src="/assets/js/page-back.js?v=1"></script>`})
		if strings.Index(body, `class="page-breadcrumb"`) > strings.Index(body, `<article class="article-shell`) {
			t.Errorf("%s: the return link must sit above the article", path)
		}
		link := regexp.MustCompile(`(?s)<a class="page-back".*?</a>`).FindString(body)
		if text := strings.TrimSpace(regexp.MustCompile(`(?s)<[^>]*>`).ReplaceAllString(link, "")); text != "" {
			t.Errorf("%s: the return link must be the arrow alone, got text %q", path, text)
		}
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/post/hello-world", nil))
	if strings.Contains(w.Body.String(), `class="page-back"`) {
		t.Error("blog posts must not carry a return arrow above the body")
	}
}
