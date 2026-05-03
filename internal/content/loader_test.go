package content

import (
	"os"
	"path/filepath"
	"strings"
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

func TestLoadPosts_SupportsNeorgFrontMatter(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "post.norg", `@document.meta
title: Norg Post
slug: norg-post
date: 2026-05-01
summary: from neorg frontmatter
tags:
  - norg
  - notes
draft: false
@end
# Norg heading
Body.
`)

	posts, err := LoadPosts(dir)
	if err != nil {
		t.Fatalf("LoadPosts returned error: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	if posts[0].Slug != "norg-post" {
		t.Fatalf("expected slug norg-post, got %q", posts[0].Slug)
	}

	if len(posts[0].Tags) != 2 || posts[0].Tags[0] != "norg" || posts[0].Tags[1] != "notes" {
		t.Fatalf("expected tags [norg notes], got %#v", posts[0].Tags)
	}
}

func TestLoadPosts_SupportsNeorgTaskRendering(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "task.norg", `@document.meta
title: Tasks
slug: tasks
date: 2026-05-01
@end
*** TODO first task
`)

	posts, err := LoadPosts(dir)
	if err != nil {
		t.Fatalf("LoadPosts error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if !strings.Contains(string(posts[0].BodyHTML), `data-task-state="todo"`) {
		t.Fatalf("expected checklist html, got %s", posts[0].BodyHTML)
	}
}

func TestLoadProjects_SortsNewestFirst(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "one.norg", `@document.meta
title: One
slug: one
date: 2026-01-01
summary: one
@end
body`)
	writePost(t, dir, "two.norg", `@document.meta
title: Two
slug: two
date: 2026-02-01
summary: two
@end
body`)

	projects, err := LoadProjects(dir)
	if err != nil {
		t.Fatalf("LoadProjects returned error: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	if projects[0].Slug != "two" {
		t.Fatalf("expected newest project first, got %q", projects[0].Slug)
	}
}

func TestLoadPage_LoadsSinglePage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "about.norg")
	writePost(t, dir, "about.norg", `@document.meta
title: About
slug: about
date: 2026-05-01
summary: about me
@end
* About`)

	page, err := LoadPage(path)
	if err != nil {
		t.Fatalf("LoadPage returned error: %v", err)
	}
	if page.Slug != "about" {
		t.Fatalf("expected slug about, got %q", page.Slug)
	}
}

func TestStore_BySlugAndByTag(t *testing.T) {
	store := NewStore([]Post{
		{Published: Published{Title: "One", Slug: "one", Date: "2026-01-01", Tags: []string{"go", "personal"}}},
		{Published: Published{Title: "Two", Slug: "two", Date: "2026-02-01", Tags: []string{"go"}}},
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
