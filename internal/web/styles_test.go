package web

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestSearchScriptSafeAsyncMessagesAndKeyboard(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/search.js")
	if err != nil {
		t.Fatalf("read search.js: %v", err)
	}
	js := string(jsBytes)
	if strings.Contains(js, "innerHTML") {
		t.Fatal("search.js should render API results with DOM text nodes, not innerHTML")
	}
	for _, marker := range []string{"AbortController", "signal:", "Loading…", "No results", "Search unavailable", "textContent", "trigger.focus()", "ctrlKey", "metaKey", `key === "escape"`} {
		if !strings.Contains(js, marker) {
			t.Errorf("search.js missing %q", marker)
		}
	}
	if strings.Contains(js, `resultLink("#")`) {
		t.Error("search status messages must not be links")
	}
}

func TestSearchResultLinksRejectUnsafePaths(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/search.js")
	if err != nil {
		t.Fatalf("read search.js: %v", err)
	}
	js := string(jsBytes)
	for _, guard := range []string{`/^\/(?!\/)/`, `!href.includes("\\")`} {
		if !strings.Contains(js, guard) {
			t.Errorf("search result URL guard missing %q", guard)
		}
	}
}

func TestTailwindRetainsDynamicSearchStyles(t *testing.T) {
	configBytes, err := os.ReadFile("../../tailwind.config.js")
	if err != nil {
		t.Fatalf("read tailwind config: %v", err)
	}
	if !strings.Contains(string(configBytes), `./internal/web/assets/js/**/*.js`) {
		t.Error("tailwind content must scan JavaScript assets")
	}

	cssBytes, err := os.ReadFile("assets/css/output.css")
	if err != nil {
		t.Fatalf("read generated CSS: %v", err)
	}
	css := string(cssBytes)
	for _, selector := range []string{
		".search-result",
		".search-result-message",
		".search-result-message-loading",
		".search-result-message-error",
		".search-result-title",
	} {
		if !strings.Contains(css, selector) {
			t.Errorf("generated CSS missing dynamic selector %q", selector)
		}
	}
}

func TestThemeScriptKeepsModesStorageAndCompactText(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/theme-toggle.js")
	if err != nil {
		t.Fatalf("read theme-toggle.js: %v", err)
	}
	js := string(jsBytes)
	for _, marker := range []string{`"theme-preference"`, `"system", "light", "dark"`, "matchMedia", "localStorage", "button.textContent = labels[mode]", `button.setAttribute("aria-label", label)`} {
		if !strings.Contains(js, marker) {
			t.Errorf("theme-toggle.js missing %q", marker)
		}
	}
}

func TestThemeRenderedShellStartsCompactWithFullLabel(t *testing.T) {
	var output bytes.Buffer
	if err := Base("Title", "Description").Render(context.Background(), &output); err != nil {
		t.Fatalf("render base: %v", err)
	}
	html := output.String()
	if !strings.Contains(html, `id="theme-toggle"`) || !strings.Contains(html, `aria-label="Theme: System"`) || !strings.Contains(html, `>System</button>`) {
		t.Fatalf("theme control should render compact text and full accessible label: %s", html)
	}
	if !strings.Contains(html, `id="search-results" class="search-results" aria-live="polite"`) {
		t.Error("search results should announce updates")
	}
}

func TestTypographyAssetsAndTokens(t *testing.T) {
	fontFiles := []string{
		"assets/fonts/source-serif-4-regular.woff2",
		"assets/fonts/source-serif-4-semibold.woff2",
		"assets/fonts/source-sans-3-regular.woff2",
		"assets/fonts/source-sans-3-semibold.woff2",
	}
	for _, path := range fontFiles {
		font, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read typography asset %q: %v", path, err)
			continue
		}
		if len(font) < 4 || string(font[:4]) != "wOF2" {
			t.Errorf("typography asset %q does not have WOFF2 magic", path)
		}
	}
	for _, path := range []string{"assets/fonts/OFL.txt", "assets/fonts/README.md"} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected typography documentation %q: %v", path, err)
		}
	}

	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	markers := []string{
		"Source Serif 4", "Source Sans 3", "--font-editorial", "--font-ui", "--font-code",
		"font-display: swap", "font-weight: 400", "font-weight: 600",
		`url("/assets/fonts/source-serif-4-regular.woff2")`,
		`url("/assets/fonts/source-serif-4-semibold.woff2")`,
		`url("/assets/fonts/source-sans-3-regular.woff2")`,
		`url("/assets/fonts/source-sans-3-semibold.woff2")`,
	}
	for _, marker := range markers {
		if !strings.Contains(string(cssBytes), marker) {
			t.Errorf("expected input.css to contain %q", marker)
		}
	}
}

func TestFlatFullWidthMasthead(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	rule := regexp.MustCompile(`(?s)\.site-masthead\s*\{[^}]*\}`).FindString(string(cssBytes))
	for _, marker := range []string{"width: 100%", "margin: 0", "border: 0", "border-radius: 0", "box-shadow: none"} {
		if !strings.Contains(rule, marker) {
			t.Errorf("flat masthead rule missing %q: %s", marker, rule)
		}
	}
}

func TestProjectContentStaysContained(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	selectors := map[string]string{
		"article shell":        `\.article-shell,\s*\.post-prose`,
		"article body":         `\.article-body`,
		"section block":        `\.section-block`,
		"project subpost grid": `\.project-subposts-grid`,
	}
	for name, selector := range selectors {
		rule := regexp.MustCompile(`(?s)` + selector + `\s*\{[^}]*\}`).FindString(css)
		if !strings.Contains(rule, "min-width: 0") {
			t.Errorf("%s must contain its content", name)
		}
	}

	shellRule := regexp.MustCompile(`(?s)\.article-shell,\s*\.post-prose\s*\{[^}]*\}`).FindString(css)
	bodyRule := regexp.MustCompile(`(?s)\.article-body\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(shellRule, "width: 100%") || strings.Contains(shellRule, "--content-reading") {
		t.Errorf("article shell must use full content width: %s", shellRule)
	}
	if !strings.Contains(bodyRule, "max-width: none") {
		t.Errorf("article body must not cap reading width: %s", bodyRule)
	}
}

func TestProjectSubpostGrid(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)
	for _, forbidden := range []string{"[data-count=", "nth-child(6n+1)", "nth-child(6n + 1)", "nth-child(6n+"} {
		if strings.Contains(css, forbidden) {
			t.Errorf("count-specific project grid CSS forbidden: %q", forbidden)
		}
	}
	for _, marker := range []string{".project-subposts-grid", "grid-template-columns: minmax(0, 1fr)", "repeat(2, minmax(0, 1fr))", "repeat(3, minmax(0, 1fr))", "grid-column: span 2", "grid-row: span 2", ".project-subpost-card", ".project-subpost-link", "min-height: 44px", ":focus-within"} {
		if !strings.Contains(css, marker) {
			t.Errorf("expected bento project grid CSS %q", marker)
		}
	}
}

func TestProjectSubpostLongContentWraps(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	for selector, property := range map[string]string{
		`.project-subpost-link`:    `min-width\s*:\s*0`,
		`.project-subpost-title`:   `overflow-wrap\s*:\s*anywhere`,
		`.project-subpost-summary`: `overflow-wrap\s*:\s*anywhere`,
	} {
		rule := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(selector) + `\s*\{[^}]*\}`).Find(cssBytes)
		if rule == nil {
			t.Errorf("%s rule missing", selector)
			continue
		}
		if !regexp.MustCompile(property).Match(rule) {
			t.Errorf("%s must contain %s", selector, property)
		}
	}
}

func TestProjectSubpostDatePreservesMutedContrast(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	dateRule := regexp.MustCompile(`(?s)\.project-subpost-date\s*\{[^}]*\}`).Find(cssBytes)
	if dateRule == nil {
		t.Fatal("project subpost date rule missing")
	}
	if strings.Contains(string(dateRule), "var(--muted)") && regexp.MustCompile(`\bopacity\s*:`).Match(dateRule) {
		t.Fatal("project subpost date must not reduce muted text contrast with opacity")
	}
}

func TestTaskLabelsUseThemeContrastTokens(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)
	for _, token := range []string{"--task-todo:", "--task-doing:", "--task-done:", "--task-cancelled:"} {
		if strings.Count(css, token) < 3 {
			t.Errorf("expected %s in light, dark, and system-dark themes", token)
		}
		if !strings.Contains(css, "color: rgb(var("+strings.TrimSuffix(token, ":")+"))") {
			t.Errorf("expected task label to use %s", token)
		}
	}
	taskStart := strings.Index(css, ".post-prose ul.task-list")
	taskEnd := strings.Index(css[taskStart:], ".search-modal")
	taskCSS := css[taskStart : taskStart+taskEnd]
	for _, fixed := range []string{"rgb(234 88 12)", "rgb(37 99 235)", "rgb(22 163 74)", "rgb(239 68 68)"} {
		if strings.Contains(taskCSS, fixed) {
			t.Errorf("fixed task label color forbidden: %s", fixed)
		}
	}
	cancelledRule := regexp.MustCompile(`(?s)li\[data-task-state="cancelled"\]\s*\{[^}]*\}`).FindString(css)
	if strings.Contains(cancelledRule, "opacity:") {
		t.Error("cancelled task must not reduce contrast with opacity")
	}
}

func TestTableFitsContentAndKeepsOverflow(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	rule := regexp.MustCompile(`(?s)\.post-prose table\s*\{[^}]*\}`).FindString(string(cssBytes))
	for _, marker := range []string{"width: fit-content", "max-width: 100%", "overflow-x: auto"} {
		if !strings.Contains(rule, marker) {
			t.Errorf("table rule missing %q: %s", marker, rule)
		}
	}
}

func TestSyntaxTokensAndProseWrapping(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	for _, token := range []string{"--syntax-keyword", "--syntax-function", "--syntax-string", "--syntax-number"} {
		if strings.Count(css, token+":") < 3 {
			t.Errorf("expected %s in light, dark, and system-dark themes", token)
		}
		if !strings.Contains(css, "color: rgb(var("+token+"))") {
			t.Errorf("expected syntax rules to use %s", token)
		}
	}
	for _, fixed := range []string{"rgb(239 68 68)", "rgb(168 85 247)", "rgb(34 197 94)", "rgb(56 189 248)"} {
		if strings.Contains(css, fixed) {
			t.Errorf("fixed syntax color forbidden: %s", fixed)
		}
	}

	articleRule := regexp.MustCompile(`(?s)\.article-body\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(articleRule, "overflow-wrap: anywhere") {
		t.Error("article body must wrap long prose and inline code")
	}
	preRule := regexp.MustCompile(`(?s)\.post-prose pre\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(preRule, "overflow-x: auto") || !strings.Contains(preRule, "overflow-wrap: normal") {
		t.Error("preformatted code must preserve horizontal scrolling without wrapping")
	}
}

func TestInputStylesUseSingleComponentsLayer(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	if got := strings.Count(string(cssBytes), "@layer components"); got != 1 {
		t.Fatalf("expected one components layer, got %d", got)
	}
}

func TestInputStyles(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}

	css := string(cssBytes)
	required := []string{
		"--bg:",
		"--surface:",
		"--surface-raised:",
		"--text:",
		"--muted:",
		"--border:",
		"--accent:",
		"--accent-strong:",
		"--focus:",
		"--shadow-color:",
		"--radius-sm:",
		"--radius-md:",
		"--content-wide:",
		"--content-reading:",
		":root[data-theme=\"dark\"]",
		":root[data-theme=\"system\"]",
		"@media (prefers-color-scheme: dark)",
		":focus-visible",
		"::selection",
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
