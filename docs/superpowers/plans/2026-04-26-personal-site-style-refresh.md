# Personal Site Soft Editorial Style Refresh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Apply a moderate antfu-inspired “Soft Editorial” visual refresh to home/blog/post pages, with light/dark/system theme support and a manual header toggle.

**Architecture:** Keep rendering server-side with templ, add stable style hooks in templates, and centralize visual behavior in `internal/web/styles/input.css`. Implement theme preference with two tiny static JS files (`theme-init.js` and `theme-toggle.js`) served from embedded assets.

**Tech Stack:** Go `net/http`, templ, Tailwind CSS, embedded static assets (`embed.FS`), Go test.

---

## File Structure Map

### Create
- `internal/web/assets/js/theme-init.js` — early theme preference hydration from `localStorage`.
- `internal/web/assets/js/theme-toggle.js` — cycles `system -> light -> dark`, updates button text/ARIA, persists preference.
- `internal/web/styles_test.go` — stylesheet guardrail test for required selectors/tokens.

### Modify
- `internal/http/blog_routes_test.go` — add failing/passing style/theme integration checks via HTTP responses.
- `internal/web/base.templ` — new shell classes, theme toggle button, theme scripts.
- `internal/web/pages.templ` — post list/detail semantic class hooks and refined structure.
- `internal/web/styles/input.css` — theme tokens, component styles, prose tuning, reduced-motion support.
- `internal/web/base_templ.go` (generated)
- `internal/web/pages_templ.go` (generated)
- `internal/web/assets/css/output.css` (generated)

---

### Task 1: Add theme toggle scaffolding (tests first)

**Files:**
- Modify: `internal/http/blog_routes_test.go`
- Modify: `internal/web/base.templ`
- Create: `internal/web/assets/js/theme-init.js`
- Create: `internal/web/assets/js/theme-toggle.js`
- Generated: `internal/web/base_templ.go`

- [ ] **Step 1: Add failing HTTP tests for theme scaffolding and theme assets**

Append these test functions to `internal/http/blog_routes_test.go`:

```go
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
		path       string
		contains   string
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

func assertContainsAll(t *testing.T, body string, markers []string) {
	t.Helper()
	for _, marker := range markers {
		if !strings.Contains(body, marker) {
			t.Fatalf("expected body to contain %q", marker)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify RED state**

Run:

```bash
go test ./internal/http -run 'TestBaseTemplateIncludesThemeControls|TestThemeAssetsAreServed' -v
```

Expected:
- FAIL because base template does not include theme scripts/toggle markers yet.
- FAIL because `/assets/js/theme-init.js` and `/assets/js/theme-toggle.js` do not exist yet.

- [ ] **Step 3: Implement minimal theme shell and JS assets**

Replace `internal/web/base.templ` with:

```templ
package web

templ Base(pageTitle string, pageDescription string) {
	<!DOCTYPE html>
	<html lang="en" class="h-full" data-theme="system">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width,initial-scale=1"/>
			<title>{ pageTitle }</title>
			<meta name="description" content={ pageDescription }/>
			<link href="/assets/css/output.css" rel="stylesheet"/>
			<link rel="alternate" type="application/rss+xml" href="/rss.xml" title="RSS"/>
			<script src="/assets/js/theme-init.js"></script>
			<script defer src="/assets/js/theme-toggle.js"></script>
		</head>
		<body class="site-shell">
			<main class="site-container">
				<header class="site-header">
					<nav class="site-nav" aria-label="Primary">
						<a href="/" class="site-nav-link">Home</a>
						<a href="/blog" class="site-nav-link">Blog</a>
						<a href="/rss.xml" class="site-nav-link">RSS</a>
					</nav>
					<button id="theme-toggle" class="theme-toggle" type="button" aria-label="Theme: System" aria-live="polite">
						Theme: System
					</button>
				</header>
				{ children... }
			</main>
		</body>
	</html>
}
```

Create `internal/web/assets/js/theme-init.js`:

```javascript
(() => {
  const key = "theme-preference";
  const validModes = new Set(["system", "light", "dark"]);

  let mode = "system";
  try {
    const stored = localStorage.getItem(key);
    if (stored && validModes.has(stored)) mode = stored;
  } catch {
    mode = "system";
  }

  document.documentElement.setAttribute("data-theme", mode);
})();
```

Create `internal/web/assets/js/theme-toggle.js`:

```javascript
(() => {
  const key = "theme-preference";
  const modes = ["system", "light", "dark"];
  const labels = {
    system: "System",
    light: "Light",
    dark: "Dark",
  };

  const root = document.documentElement;
  const button = document.getElementById("theme-toggle");
  if (!button) return;

  const normalize = (value) => (modes.includes(value) ? value : "system");

  const readPreference = () => {
    try {
      return normalize(localStorage.getItem(key));
    } catch {
      return normalize(root.getAttribute("data-theme"));
    }
  };

  const writePreference = (value) => {
    try {
      localStorage.setItem(key, value);
    } catch {
      // Ignore storage failures (private mode, etc.)
    }
  };

  const apply = (value) => {
    const mode = normalize(value);
    root.setAttribute("data-theme", mode);
    const label = `Theme: ${labels[mode]}`;
    button.textContent = label;
    button.setAttribute("aria-label", label);
    button.dataset.themeMode = mode;
  };

  const cycle = (value) => {
    const idx = modes.indexOf(normalize(value));
    return modes[(idx + 1) % modes.length];
  };

  apply(readPreference());

  button.addEventListener("click", () => {
    const next = cycle(readPreference());
    writePreference(next);
    apply(next);
  });

  const media = window.matchMedia("(prefers-color-scheme: dark)");
  const handleSchemeChange = () => {
    if (readPreference() === "system") apply("system");
  };

  if (typeof media.addEventListener === "function") {
    media.addEventListener("change", handleSchemeChange);
  } else if (typeof media.addListener === "function") {
    media.addListener(handleSchemeChange);
  }
})();
```

- [ ] **Step 4: Generate templ and verify GREEN for Task 1 tests**

Run:

```bash
templ generate -path .
go test ./internal/http -run 'TestBaseTemplateIncludesThemeControls|TestThemeAssetsAreServed' -v
```

Expected:
- PASS for both tests.

- [ ] **Step 5: Commit Task 1**

```bash
git add internal/http/blog_routes_test.go internal/web/base.templ internal/web/base_templ.go internal/web/assets/js/theme-init.js internal/web/assets/js/theme-toggle.js
git commit -m "feat: add base theme toggle scaffolding and static theme assets"
```

---

### Task 2: Add semantic page style hooks (tests first)

**Files:**
- Modify: `internal/http/blog_routes_test.go`
- Modify: `internal/web/pages.templ`
- Generated: `internal/web/pages_templ.go`

- [ ] **Step 1: Add failing tests for page-level style hooks**

Append this function to `internal/http/blog_routes_test.go`:

```go
func TestPagesExposeStyleHooks(t *testing.T) {
	s := testBlogServer()
	h := s.RegisterRoutes()

	cases := []struct {
		path    string
		markers []string
	}{
		{
			path: "/",
			markers: []string{"page-title", "page-subtitle", "section-title", "post-list", "post-item", "post-link", "tag-chip"},
		},
		{
			path: "/blog",
			markers: []string{"page-title", "post-list", "post-meta-row"},
		},
		{
			path: "/post/hello-world",
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
```

- [ ] **Step 2: Run tests to verify RED state**

Run:

```bash
go test ./internal/http -run TestPagesExposeStyleHooks -v
```

Expected:
- FAIL because the semantic class hooks are not present in `pages.templ` yet.

- [ ] **Step 3: Implement page template hook classes**

Replace `internal/web/pages.templ` with:

```templ
package web

import "github.com/danielscoffee/danielscoffee.me/internal/content"

templ HomePage(posts []content.Post) {
	@Base("Daniel's Site", "Personal website and blog") {
		<section class="intro-block">
			<h1 class="page-title">Daniel's Site</h1>
			<p class="page-subtitle">Personal notes on software, projects, and experiments.</p>
		</section>
		<section class="section-block">
			<h2 class="section-title">Recent posts</h2>
			@PostList(posts)
		</section>
	}
}

templ BlogIndexPage(posts []content.Post) {
	@Base("Blog", "All blog posts") {
		<section class="intro-block">
			<h1 class="page-title">Blog</h1>
			<p class="page-subtitle">All writing in one place.</p>
		</section>
		<section class="section-block">
			<h2 class="section-title">Latest entries</h2>
			@PostList(posts)
		</section>
	}
}

templ TagPage(tag string, posts []content.Post) {
	@Base("Tag: "+tag, "Posts filtered by tag") {
		<section class="intro-block">
			<h1 class="page-title">Tagged with <span class="tag-emphasis">{ tag }</span></h1>
			<p class="page-subtitle">Filtered writing by topic.</p>
		</section>
		<section class="section-block">
			if len(posts) == 0 {
				<p class="empty-state">No posts found for this tag yet.</p>
			} else {
				@PostList(posts)
			}
		</section>
	}
}

templ PostList(posts []content.Post) {
	<ul class="post-list">
		for _, post := range posts {
			<li class="post-item">
				<a class="post-link" href={ templ.SafeURL("/post/" + post.Slug) }>{ post.Title }</a>
				<p class="post-summary">{ post.Summary }</p>
				<div class="post-meta-row">
					<p class="post-date">{ post.Date }</p>
					if len(post.Tags) > 0 {
						<ul class="tag-list">
							for _, tag := range post.Tags {
								<li>
									<a class="tag-chip" href={ templ.SafeURL("/tag/" + tag) }>{ tag }</a>
								</li>
							}
						</ul>
					}
				</div>
			</li>
		}
	</ul>
}

templ BlogPostPage(post content.Post) {
	@Base(post.Title, post.Summary) {
		<article class="post-prose">
			<header class="post-header">
				<h1 class="post-title">{ post.Title }</h1>
				<p class="post-date">{ post.Date }</p>
			</header>
			@templ.Raw(string(post.BodyHTML))
		</article>
	}
}
```

- [ ] **Step 4: Generate templ and verify GREEN for Task 2 test**

Run:

```bash
templ generate -path .
go test ./internal/http -run TestPagesExposeStyleHooks -v
```

Expected:
- PASS for `TestPagesExposeStyleHooks`.

- [ ] **Step 5: Commit Task 2**

```bash
git add internal/http/blog_routes_test.go internal/web/pages.templ internal/web/pages_templ.go
git commit -m "feat: add semantic style hooks to home blog and post templates"
```

---

### Task 3: Implement theme tokens and soft-editorial component styling (tests first)

**Files:**
- Create: `internal/web/styles_test.go`
- Modify: `internal/web/styles/input.css`
- Generated: `internal/web/assets/css/output.css`

- [ ] **Step 1: Add failing stylesheet guardrail test**

Create `internal/web/styles_test.go`:

```go
package web

import (
	"os"
	"strings"
	"testing"
)

func TestInputStylesDefineThemeAndComponentSelectors(t *testing.T) {
	cssBytes, err := os.ReadFile("internal/web/styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	css := string(cssBytes)
	required := []string{
		":root[data-theme=\"dark\"]",
		":root[data-theme=\"system\"]",
		".site-nav-link",
		".theme-toggle",
		".post-list",
		".post-item",
		".post-prose",
		".tag-chip",
		"@media (prefers-reduced-motion: reduce)",
	}

	for _, marker := range required {
		if !strings.Contains(css, marker) {
			t.Fatalf("expected input.css to contain %q", marker)
		}
	}
}
```

- [ ] **Step 2: Run test to verify RED state**

Run:

```bash
go test ./internal/web -run TestInputStylesDefineThemeAndComponentSelectors -v
```

Expected:
- FAIL because current `input.css` only has Tailwind directives.

- [ ] **Step 3: Implement full soft-editorial stylesheet**

Replace `internal/web/styles/input.css` with:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    color-scheme: light;
    --bg: 250 250 249;
    --surface: 255 255 255;
    --text: 24 24 27;
    --muted: 82 82 91;
    --border: 228 228 231;
    --accent: 63 63 70;
    --accent-strong: 24 24 27;
    --chip-bg: 244 244 245;
    --chip-text: 63 63 70;
  }

  :root[data-theme="dark"] {
    color-scheme: dark;
    --bg: 10 10 12;
    --surface: 24 24 27;
    --text: 228 228 231;
    --muted: 161 161 170;
    --border: 63 63 70;
    --accent: 212 212 216;
    --accent-strong: 244 244 245;
    --chip-bg: 39 39 42;
    --chip-text: 212 212 216;
  }

  @media (prefers-color-scheme: dark) {
    :root[data-theme="system"] {
      color-scheme: dark;
      --bg: 10 10 12;
      --surface: 24 24 27;
      --text: 228 228 231;
      --muted: 161 161 170;
      --border: 63 63 70;
      --accent: 212 212 216;
      --accent-strong: 244 244 245;
      --chip-bg: 39 39 42;
      --chip-text: 212 212 216;
    }
  }

  html,
  body {
    @apply h-full;
  }

  body {
    background-color: rgb(var(--bg));
    color: rgb(var(--text));
    font-feature-settings: "rlig" 1, "calt" 1;
  }

  a {
    transition: color 180ms ease, border-color 180ms ease, background-color 180ms ease;
  }
}

@layer components {
  .site-shell {
    @apply min-h-full;
  }

  .site-container {
    @apply mx-auto max-w-3xl px-4 py-10 md:py-12;
  }

  .site-header {
    @apply mb-10 flex items-center justify-between gap-4;
  }

  .site-nav {
    @apply flex items-center gap-4 text-sm;
  }

  .site-nav-link {
    @apply border-b border-transparent pb-0.5 no-underline;
    color: rgb(var(--muted));
  }

  .site-nav-link:hover,
  .site-nav-link:focus-visible {
    color: rgb(var(--accent-strong));
    border-color: rgb(var(--border));
  }

  .theme-toggle {
    @apply rounded-md border px-2.5 py-1 text-xs font-medium tracking-wide;
    border-color: rgb(var(--border));
    color: rgb(var(--muted));
    background-color: rgb(var(--surface));
  }

  .theme-toggle:hover,
  .theme-toggle:focus-visible {
    color: rgb(var(--accent-strong));
    border-color: rgb(var(--accent));
  }

  .intro-block {
    @apply space-y-3;
  }

  .section-block {
    @apply mt-10 space-y-4;
  }

  .page-title {
    @apply text-4xl font-semibold tracking-tight;
  }

  .section-title {
    @apply text-2xl font-semibold tracking-tight;
  }

  .page-subtitle {
    @apply text-base leading-7;
    color: rgb(var(--muted));
  }

  .tag-emphasis {
    color: rgb(var(--accent-strong));
  }

  .empty-state {
    @apply text-sm;
    color: rgb(var(--muted));
  }

  .post-list {
    @apply mt-4 space-y-6;
  }

  .post-item {
    @apply border-b pb-6;
    border-color: rgb(var(--border));
  }

  .post-link {
    @apply text-2xl font-semibold no-underline;
    color: rgb(var(--accent-strong));
    border-bottom: 1px solid transparent;
  }

  .post-link:hover,
  .post-link:focus-visible {
    border-bottom-color: rgb(var(--border));
    color: rgb(var(--text));
  }

  .post-summary {
    @apply mt-2 leading-7;
    color: rgb(var(--muted));
  }

  .post-meta-row {
    @apply mt-3 flex flex-wrap items-center gap-3;
  }

  .post-date {
    @apply text-sm;
    color: rgb(var(--muted));
  }

  .tag-list {
    @apply flex flex-wrap gap-2;
  }

  .tag-chip {
    @apply rounded-md px-2 py-1 text-xs no-underline;
    background-color: rgb(var(--chip-bg));
    color: rgb(var(--chip-text));
  }

  .tag-chip:hover,
  .tag-chip:focus-visible {
    color: rgb(var(--accent-strong));
  }

  .post-prose {
    @apply mt-8 text-[1.03rem] leading-8;
    color: rgb(var(--text));
  }

  .post-header {
    @apply mb-8 space-y-2 border-b pb-5;
    border-color: rgb(var(--border));
  }

  .post-title {
    @apply text-4xl font-semibold tracking-tight;
  }

  .post-prose h2,
  .post-prose h3,
  .post-prose h4 {
    @apply mt-10 mb-3 font-semibold tracking-tight;
  }

  .post-prose p,
  .post-prose ul,
  .post-prose ol {
    @apply my-4;
    color: rgb(var(--text));
  }

  .post-prose a {
    color: rgb(var(--accent-strong));
    text-decoration: underline;
    text-decoration-color: rgb(var(--border));
    text-underline-offset: 3px;
  }

  .post-prose blockquote {
    @apply my-6 border-l-2 pl-4 italic;
    border-color: rgb(var(--border));
    color: rgb(var(--muted));
  }

  .post-prose code {
    @apply rounded px-1.5 py-0.5 text-[0.92em];
    background-color: rgb(var(--chip-bg));
  }

  .post-prose pre {
    @apply my-6 overflow-x-auto rounded-md border p-4 text-sm leading-6;
    border-color: rgb(var(--border));
    background-color: rgb(var(--surface));
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

- [ ] **Step 4: Verify GREEN and regenerate Tailwind output**

Run:

```bash
go test ./internal/web -run TestInputStylesDefineThemeAndComponentSelectors -v
make generate
```

Expected:
- Stylesheet test PASS.
- `make generate` succeeds and updates `internal/web/assets/css/output.css`.

- [ ] **Step 5: Commit Task 3**

```bash
git add internal/web/styles_test.go internal/web/styles/input.css internal/web/assets/css/output.css
git commit -m "feat: add soft-editorial theme tokens and component styles"
```

---

## Final Verification Gate

- [ ] **Step 1: Run full automated verification**

```bash
templ generate -path .
go test ./...
go build ./...
```

Expected:
- All commands succeed with zero failing tests.

- [ ] **Step 2: Run manual visual verification**

Run app:

```bash
make run
```

Manual checks in browser (`http://localhost:8080`):
1. `/` uses refreshed spacing/typography and subtle metadata styling.
2. `/blog` shows soft editorial list styling with gentle interactions.
3. `/post/hello-world` renders refined prose, headers, code, and blockquotes.
4. Theme toggle cycles `System -> Light -> Dark -> System` and persists after reload.
5. With reduced motion enabled at OS level, transitions are effectively minimized.

Stop app after checks.

---

## Self-Review (plan quality)

- **Spec coverage:** all approved sections covered: visual language, component updates, theme toggle behavior, reduced motion, and verification.
- **Placeholder scan:** no TODO/TBD placeholders remain.
- **Type/signature consistency:** all test names, selectors, and file paths are consistent across tasks.
