package web

import (
	"os"
	"strings"
	"testing"
)

func TestSearchScriptAvoidsInnerHTMLForResults(t *testing.T) {
	jsBytes, err := os.ReadFile("assets/js/search.js")
	if err != nil {
		t.Fatalf("read search.js: %v", err)
	}

	if strings.Contains(string(jsBytes), "innerHTML") {
		t.Fatalf("search.js should render API results with DOM text nodes, not innerHTML")
	}
}

func TestTypographyAssetsAndTokens(t *testing.T) {
	fontFiles := []string{
		"assets/fonts/source-serif-4-regular.woff2",
		"assets/fonts/source-serif-4-semibold.woff2",
		"assets/fonts/source-sans-3-regular.woff2",
		"assets/fonts/source-sans-3-semibold.woff2",
		"assets/fonts/OFL.txt",
	}
	for _, path := range fontFiles {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected typography asset %q: %v", path, err)
		}
	}

	cssBytes, err := os.ReadFile("styles/input.css")
	if err != nil {
		t.Fatalf("read styles: %v", err)
	}
	for _, marker := range []string{"Source Serif 4", "Source Sans 3", "--font-editorial", "--font-ui", "--font-code"} {
		if !strings.Contains(string(cssBytes), marker) {
			t.Errorf("expected input.css to contain %q", marker)
		}
	}
}

func TestInputStylesDefineThemeAndComponentSelectors(t *testing.T) {
	cssBytes, err := os.ReadFile("styles/input.css")
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
