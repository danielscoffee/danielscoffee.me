package httpapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

func TestSearchRejectsOversizedQuery(t *testing.T) {
	h := testBlogServer().RegisterRoutes()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/search?q="+strings.Repeat("x", 257), nil)

	h.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected JSON response, got %q", got)
	}
}

func TestSearchIndexer_FilterAndRanking(t *testing.T) {
	docs := []content.SearchDoc{
		{Type: "blog", Title: "Go internals", Slug: "go-internals", Date: "2026-05-01", Summary: "summary", Body: "details"},
		{Type: "blog", Title: "Notes", Slug: "notes", Date: "2026-05-02", Summary: "go language", Body: "go in body"},
		{Type: "projects", Title: "Site project", Slug: "site-project", Date: "2026-05-03", Summary: "project summary", Body: "go tooling"},
	}

	idx := NewSearchIndexer(docs)

	projectResults := idx.Search("projects go")
	if len(projectResults) != 1 || projectResults[0].Type != "projects" {
		t.Fatalf("expected only projects result, got %#v", projectResults)
	}

	blogResults := idx.Search("blog go")
	if len(blogResults) < 2 {
		t.Fatalf("expected at least 2 blog results, got %#v", blogResults)
	}
	if blogResults[0].Slug != "go-internals" {
		t.Fatalf("expected title match ranked first, got %#v", blogResults)
	}
}

func TestSearchIndexer_SubPostURL(t *testing.T) {
	docs := []content.SearchDoc{
		{Type: "projects", Title: "Rebuild", Slug: "rebuild", ParentSlug: "blog-project", Date: "2026-05-10", Summary: "rebuild"},
	}
	idx := NewSearchIndexer(docs)
	results := idx.Search("projects rebuild")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].URL != "/projects/blog-project/rebuild" {
		t.Fatalf("expected nested subpost URL, got %q", results[0].URL)
	}
}

func TestSearchIndexer_AllItemsForBareFilter(t *testing.T) {
	docs := []content.SearchDoc{
		{Type: "blog", Title: "Go internals", Slug: "go-internals", Date: "2026-05-01"},
		{Type: "blog", Title: "Notes", Slug: "notes", Date: "2026-05-02"},
		{Type: "projects", Title: "Site project", Slug: "site-project", Date: "2026-05-03"},
	}

	idx := NewSearchIndexer(docs)

	blogResults := idx.Search("blog")
	if len(blogResults) != 2 {
		t.Fatalf("expected all blog results, got %#v", blogResults)
	}
	for _, result := range blogResults {
		if result.Type != "blog" {
			t.Fatalf("expected only blog results, got %#v", blogResults)
		}
	}

	projectResults := idx.Search("projects")
	if len(projectResults) != 1 || projectResults[0].Type != "projects" {
		t.Fatalf("expected only project results, got %#v", projectResults)
	}
}
