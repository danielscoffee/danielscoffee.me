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

// parentPath is where the masthead's back arrow leads: one level up the
// site's tree, whatever page the reader came from.
func parentPath(ctx context.Context) string {
	path := strings.TrimRight(pageContext(ctx).Path, "/")
	switch {
	case strings.HasPrefix(path, "/post/"), strings.HasPrefix(path, "/tag/"):
		return "/blog"
	case strings.HasPrefix(path, "/projects/"):
		return path[:strings.LastIndex(path, "/")]
	default:
		return "/"
	}
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
