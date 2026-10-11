package norg

import (
	"fmt"
	"html"
	"strings"
)

type codeLanguage struct {
	label string
	path  string // SVG path data on a 24x24 grid
}

// codeLanguageAliases maps the names people write after @code to a key in
// codeLanguages.
var codeLanguageAliases = map[string]string{
	"golang":     "go",
	"htm":        "html",
	"js":         "javascript",
	"mjs":        "javascript",
	"ts":         "typescript",
	"sh":         "bash",
	"shell":      "bash",
	"zsh":        "bash",
	"dockerfile": "docker",
	"yml":        "yaml",
	"py":         "python",
	"rs":         "rust",
	"md":         "markdown",
	"kt":         "kotlin",
	"psql":       "sql",
	"postgres":   "sql",
	"postgresql": "sql",
}

const (
	// A drawn </> for languages without a mark, in the masthead arrow's stroke.
	genericCodeIcon = `<svg class="code-lang-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="M8 7l-5 5 5 5"></path><path d="M16 7l5 5-5 5"></path><path d="M13.5 4l-3 16"></path></svg>`
	copyIcon        = `<svg class="code-copy-icon code-copy-icon-copy" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><rect x="5.5" y="5.5" width="8" height="8" rx="1.5"></rect><path d="M10.5 5.5V4a1.5 1.5 0 0 0-1.5-1.5H4A1.5 1.5 0 0 0 2.5 4v5A1.5 1.5 0 0 0 4 10.5h1.5"></path></svg>`
	checkIcon       = `<svg class="code-copy-icon code-copy-icon-check" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="M3 8.5l3 3 7-7"></path></svg>`
)

// renderCodeBlock wraps highlighted code in a figure whose caption bar shows
// the language mark and name, plus a copy button. The button ships hidden;
// code-copy.js reveals it, so readers without JavaScript never see a dead
// control.
func renderCodeBlock(lang, code string) string {
	key, label, path := lookupCodeLanguage(lang)

	var b strings.Builder
	fmt.Fprintf(&b, `<figure class="code-block" data-lang="%s">`, html.EscapeString(key))
	b.WriteString(`<figcaption class="code-block-bar"><span class="code-lang">`)
	if path == "" {
		b.WriteString(genericCodeIcon)
	} else {
		fmt.Fprintf(&b, `<svg class="code-lang-icon" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="%s"></path></svg>`, path)
	}
	fmt.Fprintf(&b, `<span class="code-lang-name">%s</span></span>`, html.EscapeString(label))
	fmt.Fprintf(&b, `<button class="code-copy" type="button" aria-label="Copy %s code" hidden>%s%s<span class="code-copy-text">Copy</span></button>`, html.EscapeString(label), copyIcon, checkIcon)
	b.WriteString(`</figcaption>`)
	b.WriteString(renderHighlightedCode(lang, code))
	b.WriteString(`</figure>`)
	return b.String()
}

// lookupCodeLanguage returns the data-lang key, the display label, and the
// mark's path data ("" means the generic mark).
func lookupCodeLanguage(lang string) (string, string, string) {
	name := strings.ToLower(strings.TrimSpace(lang))
	if alias, ok := codeLanguageAliases[name]; ok {
		name = alias
	}
	if known, ok := codeLanguages[name]; ok {
		return name, known.label, known.path
	}
	if name == "" || name == "text" || name == "plain" {
		return "text", "Code", ""
	}
	return "other", strings.TrimSpace(lang), ""
}
