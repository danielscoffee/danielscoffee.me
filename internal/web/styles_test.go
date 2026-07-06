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
