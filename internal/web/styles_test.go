package web

import (
	"os"
	"regexp"
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
