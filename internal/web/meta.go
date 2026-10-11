package web

import (
	"context"
	"strings"
)

type pageContextKey struct{}

// PageContext carries what the page head needs to build absolute URLs: the
// public base URL from SITE_URL ("" when it is unset) and the request path.
type PageContext struct {
	BaseURL string
	Path    string
}

// WithPageContext attaches request-level page details for templates.
func WithPageContext(ctx context.Context, pc PageContext) context.Context {
	return context.WithValue(ctx, pageContextKey{}, pc)
}

func pageContext(ctx context.Context) PageContext {
	pc, _ := ctx.Value(pageContextKey{}).(PageContext)
	return pc
}

// Crumb is one link in the masthead breadcrumb below the section, such as
// the project above a devlog.
type Crumb struct {
	Label string
	Href  string
}

// shortLabel is a crumb's middle size: its first word, cut to four letters
// when longer than five ("projects" -> "proj", "Blog Engine" -> "Blog").
func shortLabel(label string) string {
	fields := strings.Fields(label)
	if len(fields) == 0 {
		return ""
	}
	word := []rune(fields[0])
	if len(word) > 5 {
		word = word[:4]
	}
	return string(word)
}

// tinyLabel is a crumb's smallest size: its first letter.
func tinyLabel(label string) string {
	word := []rune(strings.TrimSpace(label))
	if len(word) == 0 {
		return ""
	}
	return string(word[:1])
}

// isCurrentPath reports whether the request is for path itself, so its crumb
// can be marked as the current page instead of linking to itself.
func isCurrentPath(ctx context.Context, path string) bool {
	return strings.TrimRight(pageContext(ctx).Path, "/") == path
}

// absoluteURL joins the configured base URL and path. It returns "" when no
// base URL is configured, so the head never advertises a guessed origin.
func absoluteURL(ctx context.Context, path string) string {
	base := strings.TrimRight(pageContext(ctx).BaseURL, "/")
	if base == "" {
		return ""
	}
	return base + path
}

// canonicalURL is the current page without its query string, so filtered
// views such as /projects?tag=go point back at /projects.
func canonicalURL(ctx context.Context) string {
	path := pageContext(ctx).Path
	if path == "" {
		path = "/"
	}
	return absoluteURL(ctx, path)
}

// blogPostingJSONLD describes a dated post with only what the page shows.
func blogPostingJSONLD(ctx context.Context, title, date string) map[string]any {
	data := map[string]any{
		"@context": "https://schema.org",
		"@type":    "BlogPosting",
		"headline": title,
	}
	if date != "" {
		data["datePublished"] = date
	}
	if url := canonicalURL(ctx); url != "" {
		data["url"] = url
		data["mainEntityOfPage"] = url
	}
	return data
}

const (
	siteName        = "danielscoffee"
	shareImagePath  = "/assets/img/og-default.jpg"
	shareImageAlt   = "danielscoffee: Backend developer. Systems enjoyer. Production bug survivor."
	pageTitleSuffix = " · danielscoffee"
)
