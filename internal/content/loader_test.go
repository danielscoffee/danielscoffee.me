package content

import (
	"os"
	"path/filepath"
	"testing"
)

func writePost(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatalf("write post %s: %v", name, err)
	}
}

func TestLoadPosts_SortsNewestFirstAndSkipsDrafts(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "one.md", `---
title: One
slug: one
date: 2026-01-01
summary: one summary
tags: [go]
draft: false
---
# One
Body one.
`)
	writePost(t, dir, "two.md", `---
title: Two
slug: two
date: 2026-03-01
summary: two summary
tags: [personal]
draft: false
---
# Two
Body two.
`)
	writePost(t, dir, "draft.md", `---
title: Draft
slug: draft
date: 2026-04-01
summary: draft summary
tags: [draft]
draft: true
---
# Draft
`)

	posts, err := LoadPosts(dir)
	if err != nil {
		t.Fatalf("LoadPosts returned error: %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 published posts, got %d", len(posts))
	}

	if posts[0].Slug != "two" {
		t.Fatalf("expected newest post first, got slug %q", posts[0].Slug)
	}

	if posts[1].Slug != "one" {
		t.Fatalf("expected oldest post second, got slug %q", posts[1].Slug)
	}

	if posts[0].BodyHTML == "" {
		t.Fatalf("expected rendered html body for post %q", posts[0].Slug)
	}
}

func TestLoadPosts_RequiresTitleSlugDate(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "bad.md", `---
title: Missing slug
date: 2026-01-01
summary: no slug
---
# Bad
`)

	_, err := LoadPosts(dir)
	if err == nil {
		t.Fatal("expected error for missing required frontmatter")
	}
}

func TestStore_BySlugAndByTag(t *testing.T) {
	store := NewStore([]Post{
		{Title: "One", Slug: "one", Date: "2026-01-01", Tags: []string{"go", "personal"}},
		{Title: "Two", Slug: "two", Date: "2026-02-01", Tags: []string{"go"}},
	})

	if _, ok := store.BySlug("one"); !ok {
		t.Fatal("expected slug lookup for one to succeed")
	}

	if _, ok := store.BySlug("missing"); ok {
		t.Fatal("expected slug lookup for missing to fail")
	}

	posts := store.ByTag("personal")
	if len(posts) != 1 || posts[0].Slug != "one" {
		t.Fatalf("expected personal tag to return one post 'one', got %#v", posts)
	}
}
