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
	for _, marker := range []string{`"theme-preference"`, `"(prefers-color-scheme: dark)"`, `preference || "system"`, "matchMedia", "localStorage", "button.textContent = label", `button.setAttribute("aria-label", label)`} {
		if !strings.Contains(js, marker) {
			t.Errorf("theme-toggle.js missing %q", marker)
		}
	}
	if strings.Contains(js, `"System"`) {
		t.Error("theme toggle must switch light/dark only; following the OS is the default, not a button state")
	}
}

func TestThemeRenderedShellStartsCompactWithFullLabel(t *testing.T) {
	var output bytes.Buffer
	if err := Base("Title", "Description", "blog").Render(context.Background(), &output); err != nil {
		t.Fatalf("render base: %v", err)
	}
	html := output.String()
	if !strings.Contains(html, `id="theme-toggle"`) || !strings.Contains(html, `aria-label="Theme"`) || !strings.Contains(html, `>Theme</button>`) {
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

func TestShortPagesKeepFooterAtViewportEnd(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	shellRule := regexp.MustCompile(`(?s)\.site-shell\s*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"min-height: 100dvh", "display: flex", "flex-direction: column"} {
		if !strings.Contains(shellRule, marker) {
			t.Errorf("site shell missing %q: %s", marker, shellRule)
		}
	}
	containerRule := regexp.MustCompile(`(?s)\.site-container\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(containerRule, "flex: 1 0 auto") {
		t.Errorf("site content must fill short viewports: %s", containerRule)
	}
}

func TestResponsiveFooterStaysSingleLine(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	footerRules := regexp.MustCompile(`(?s)\.footer-inner\s*\{[^}]*\}`).FindAllString(css, -1)
	if !strings.Contains(strings.Join(footerRules, "\n"), "flex-direction: row") {
		t.Errorf("responsive footer must stay on one row: %v", footerRules)
	}
	footerNavRules := regexp.MustCompile(`(?s)\.footer-inner nav\s*\{[^}]*\}`).FindAllString(css, -1)
	if !strings.Contains(strings.Join(footerNavRules, "\n"), "flex-wrap: nowrap") {
		t.Errorf("responsive footer links must not wrap: %v", footerNavRules)
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

func TestBackLinkMasthead(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	sectionRule := regexp.MustCompile(`(?s)\t\.site-section\s*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"min-height: 44px", "var(--font-editorial)"} {
		if !strings.Contains(sectionRule, marker) {
			t.Errorf("section label missing %q: %s", marker, sectionRule)
		}
	}
	backRule := regexp.MustCompile(`(?s)\t\.site-back\s*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"width: 44px", "height: 44px"} {
		if !strings.Contains(backRule, marker) {
			t.Errorf("back arrow needs a 44px target, missing %q: %s", marker, backRule)
		}
	}
	for _, removed := range []string{".site-nav-desktop", ".mobile-menu", ".site-menu", ".social-link"} {
		if strings.Contains(css, removed) {
			t.Errorf("old masthead selector %q must be removed", removed)
		}
	}
}

func TestSiteBackScript(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/site-back.js")
	if err != nil {
		t.Fatalf("read site-back.js: %v", err)
	}
	js := string(jsBytes)
	for _, marker := range []string{`querySelector(".site-back")`, "document.referrer", "location.origin", "history.back()", "event.preventDefault()", "event.metaKey", "event.ctrlKey"} {
		if !strings.Contains(js, marker) {
			t.Errorf("site-back.js missing %q", marker)
		}
	}
}

func TestCodeCopyScript(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/code-copy.js")
	if err != nil {
		t.Fatalf("read code-copy.js: %v", err)
	}
	js := string(jsBytes)
	if strings.Contains(js, "innerHTML") {
		t.Error("code-copy.js must not use innerHTML")
	}
	for _, marker := range []string{"navigator.clipboard.writeText", "button.hidden = false", `"aria-live", "polite"`, "selectNodeContents", "textContent"} {
		if !strings.Contains(js, marker) {
			t.Errorf("code-copy.js missing %q", marker)
		}
	}

	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)
	if !regexp.MustCompile(`\.code-copy\[hidden\] \{\n\t\tdisplay: none;`).MatchString(css) {
		t.Error("the copy button must stay hidden until code-copy.js reveals it")
	}
	copyRule := regexp.MustCompile(`(?s)\t\.post-prose \.code-copy \{[^}]*\}`).FindString(css)
	if !strings.Contains(copyRule, "min-height: 44px") {
		t.Errorf("copy button needs a 44px target: %s", copyRule)
	}
}

func TestThemeSwitchFades(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/theme-toggle.js")
	if err != nil {
		t.Fatalf("read theme-toggle.js: %v", err)
	}
	js := string(jsBytes)
	for _, marker := range []string{
		`document.startViewTransition({ update: switchTheme, types: ["theme"] })`,
		`CSS.supports("selector(:active-view-transition-type(theme))")`,
		`root.classList.add("theme-fading")`,
		`const fadeMs = 700;`,
	} {
		if !strings.Contains(js, marker) {
			t.Errorf("theme-toggle.js missing %q", marker)
		}
	}

	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)
	curve := "cubic-bezier(0.45, 0, 0.55, 1)"
	fadeRule := regexp.MustCompile(`(?s)html:active-view-transition-type\(theme\)::view-transition-old\(root\),\nhtml:active-view-transition-type\(theme\)::view-transition-new\(root\) \{[^}]*\}`).FindString(css)
	if !strings.Contains(fadeRule, "animation-duration: 700ms") || !strings.Contains(fadeRule, curve) {
		t.Errorf("theme crossfade must run 700ms on %s: %q", curve, fadeRule)
	}
	fallback := regexp.MustCompile(`(?s)html\.theme-fading,\nhtml\.theme-fading \*,[^{]*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"transition-duration: 700ms", curve, "background-color", "color"} {
		if !strings.Contains(fallback, marker) {
			t.Errorf("theme fade fallback missing %q: %q", marker, fallback)
		}
	}
	reducedAt := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedAt < 0 || !strings.Contains(css[reducedAt:], "html.theme-fading *::after {\n\t\ttransition-duration: 700ms !important;") {
		t.Error("the theme fade moves nothing and must survive reduced motion")
	}
}

func TestBrandViewTransition(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	if !regexp.MustCompile(`(?m)^@view-transition \{\n\tnavigation: auto;`).MatchString(css) {
		t.Error("pages must opt in to cross-document view transitions")
	}
	for _, section := range []string{"blog", "projects", "about"} {
		nameRule := regexp.MustCompile(`(?s)\.door-name\[data-section="` + section + `"\],\n\.site-section\[data-section="` + section + `"\] \{[^}]*\}`).FindString(css)
		if !strings.Contains(nameRule, "view-transition-name: section-"+section) {
			t.Errorf("door name and section label for %s must share a transition name: %q", section, nameRule)
		}
	}
	reducedAt := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedAt < 0 || !strings.Contains(css[reducedAt:], "@view-transition {\n\t\tnavigation: none;") {
		t.Error("reduced motion must turn view transitions off")
	}

	cssOut, err := os.ReadFile("assets/css/output.css")
	if err != nil {
		t.Fatalf("read generated CSS: %v", err)
	}
	if !strings.Contains(string(cssOut), "@view-transition") {
		t.Error("generated CSS must keep the @view-transition rule")
	}
}

func TestHomeBoardStyles(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	// The home doors sit straight on the page (no panel), top-left, so they
	// use the page tokens and carry no board token family of their own.
	if strings.Contains(css, "--board") {
		t.Error("the home page has no board panel; --board tokens must not return")
	}
	boardRule := regexp.MustCompile(`(?s)\t\.menu-board\s*\{[^}]*\}`).FindString(css)
	if strings.Contains(boardRule, "background") || strings.Contains(boardRule, "padding") || !strings.Contains(boardRule, "color: rgb(var(--text))") {
		t.Errorf("menu board must be unboxed page text: %s", boardRule)
	}
	shellRule := regexp.MustCompile(`(?s)\t\.home-shell\s*\{[^}]*\}`).FindString(css)
	if strings.Contains(shellRule, "place-items") || strings.Contains(shellRule, "center") {
		t.Errorf("home must anchor to the top-left corner, not centre: %s", shellRule)
	}

	frameRule := regexp.MustCompile(`(?s)\.door-frame\s*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"width: 96px", "height: 96px", "overflow: hidden", "transition: width 0.3s ease-out, height 0.3s ease-out"} {
		if !strings.Contains(frameRule, marker) {
			t.Errorf("door frame missing %q: %s", marker, frameRule)
		}
	}
	imgRule := regexp.MustCompile(`(?s)\.door-frame img\s*\{[^}]*\}`).FindString(css)
	for _, marker := range []string{"width: 192px", "height: 192px", "max-width: none", "transform: translate(-50%, -50%)"} {
		if !strings.Contains(imgRule, marker) {
			t.Errorf("door scene missing %q: %s", marker, imgRule)
		}
	}

	hoverBlock := regexp.MustCompile(`(?s)@media \(hover: hover\)\s*\{.*?\n\}`).FindString(css)
	if !strings.Contains(hoverBlock, "width: 192px") || !strings.Contains(hoverBlock, ":focus-visible .door-frame") {
		t.Errorf("hover devices must grow frames to 192px on hover and focus: %s", hoverBlock)
	}
	touchBlock := regexp.MustCompile(`(?s)@media \(hover: none\)\s*\{.*?\n\}`).FindString(css)
	if !strings.Contains(touchBlock, "width: 128px") || !strings.Contains(touchBlock, "transition: none") {
		t.Errorf("touch devices must rest frames at 128px without growth: %s", touchBlock)
	}
	narrowBlock := regexp.MustCompile(`(?s)@media \(max-width: 480px\)\s*\{\n\t\.door-leader.*?\n\}`).FindString(css)
	if !strings.Contains(narrowBlock, ".door-route") || !strings.Contains(narrowBlock, "display: none") {
		t.Errorf("narrow screens must hide leader and route: %s", narrowBlock)
	}

	reducedAt := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedAt < 0 {
		t.Fatal("input.css is missing the prefers-reduced-motion block")
	}
	reduced := css[reducedAt:]
	for _, marker := range []string{".board-intro {\n\t\tdisplay: none;", ".board-row {\n\t\tanimation: none;"} {
		if !strings.Contains(reduced, marker) {
			t.Errorf("reduced motion must skip the intro: missing %q", marker)
		}
	}
	focusRule := regexp.MustCompile(`(?s)\.board-door:focus-visible\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(focusRule, "outline: 3px solid rgb(var(--focus))") {
		t.Errorf("board doors need the site focus ring: %s", focusRule)
	}
}

func TestHomePageRendersBoardOnly(t *testing.T) {
	var output bytes.Buffer
	if err := HomePage().Render(context.Background(), &output); err != nil {
		t.Fatalf("render home: %v", err)
	}
	html := output.String()
	for _, marker := range []string{`class="menu-board"`, `class="board-door" href="/blog"`, `aria-hidden="true"`, `alt=""`} {
		if !strings.Contains(html, marker) {
			t.Errorf("home page missing %q", marker)
		}
	}
	for _, absent := range []string{"site-masthead", "site-footer", "search-modal"} {
		if strings.Contains(html, absent) {
			t.Errorf("home page should not render %q", absent)
		}
	}
}

func TestHomeIntroScript(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/home-intro.js")
	if err != nil {
		t.Fatalf("read home-intro.js: %v", err)
	}
	js := string(jsBytes)
	for _, marker := range []string{"sessionStorage", `dataset.intro = "skip"`, "keydown", "pointerdown", "focusin", "catch"} {
		if !strings.Contains(js, marker) {
			t.Errorf("home-intro.js missing %q", marker)
		}
	}
	if strings.Contains(js, "innerHTML") {
		t.Error("home-intro.js must not use innerHTML")
	}
}

func TestProjectContentStaysContained(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)

	selectors := map[string]string{
		"article shell": `\.article-shell,\s*\.post-prose`,
		"article body":  `\.article-body`,
		"section block": `\.section-block`,
		"index entry":   `\.entry`,
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

func TestDevlogListIsFlat(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	css := string(cssBytes)
	for _, gone := range []string{".project-subposts-grid", ".project-subpost-card", "grid-row: span 2"} {
		if strings.Contains(css, gone) {
			t.Errorf("devlogs are a flat list now; %q must be removed", gone)
		}
	}
	linkRule := regexp.MustCompile(`(?s)\t\.entry-link \{[^}]*\}`).FindString(css)
	if !strings.Contains(linkRule, "overflow-wrap: anywhere") {
		t.Errorf("long entry titles must wrap: %s", linkRule)
	}
	metaRule := regexp.MustCompile(`(?s)\t\.entry-meta \{[^}]*\}`).FindString(css)
	if !strings.Contains(metaRule, "var(--muted)") || strings.Contains(metaRule, "opacity") {
		t.Errorf("entry meta must use the muted ink without opacity: %s", metaRule)
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
		".year-list",
		".year-mark",
		".entry-link",
		".entry-meta",
		".post-prose",
		".article-tag",
		"@media (prefers-reduced-motion: reduce)",
	}

	for _, marker := range required {
		if !strings.Contains(css, marker) {
			t.Fatalf("expected input.css to contain %q", marker)
		}
	}
}
