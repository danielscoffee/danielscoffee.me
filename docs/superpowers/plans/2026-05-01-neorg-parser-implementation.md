# Neorg Parser Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add first-class `.norg` parsing (frontmatter + body subset) so Neorg tasks like `*** TODO` render as checklist HTML with state attributes.

**Architecture:** Keep `internal/content/loader.go` as extension router. Add dedicated `internal/content/norg_parser.go` for `.norg` metadata/body parsing and HTML rendering. Keep `.md` path unchanged and fail only on Neorg structural errors.

**Tech Stack:** Go, stdlib (`strings`, `html`, `regexp`), existing `goldmark` for `.md`, `testing` package.

---

## File map

- Create: `internal/content/norg_parser.go` — Neorg metadata parser, block parser, HTML renderer.
- Create: `internal/content/norg_parser_test.go` — parser unit tests.
- Modify: `internal/content/loader.go` — route `.norg` to new parser, keep `.md` existing behavior.
- Modify: `internal/content/loader_test.go` — integration tests for `.norg` task rendering.
- Optional modify: `content/posts/neorg-post.norg` — sample content aligned to parser behavior.

### Task 1: Add failing tests for Neorg parser behavior

**Files:**
- Create: `internal/content/norg_parser_test.go`
- Modify: `internal/content/loader_test.go`

- [ ] **Step 1: Write failing unit tests for frontmatter + tasks + structural errors**

```go
package content

import "testing"

func TestParseNorg_MixedContentAndTasks(t *testing.T) {
	raw := `@document.meta
	title: Norg Post
	slug: norg-post
	date: 2026-05-01
	summary: test
	tags:
	  - norg
	  - notes
	draft: false
	@end
	* Heading
	*** TODO ship parser
	*** DOING write tests
	*** DONE pass tests
	*** CANCELLED old idea
	- bullet
	1. ordered
	` + "```go\nfmt.Println(\"ok\")\n```" + `
	[docs](https://example.com)`

	meta, body, html, err := parseNorg(raw)
	if err != nil {
		t.Fatalf("parseNorg error: %v", err)
	}
	if meta.Slug != "norg-post" {
		t.Fatalf("slug mismatch: %q", meta.Slug)
	}
	if body == "" || html == "" {
		t.Fatalf("expected non-empty body/html")
	}
	for _, s := range []string{"todo", "doing", "done", "cancelled"} {
		if !contains(html, `data-task-state="`+s+`"`) {
			t.Fatalf("missing task state %s in html: %s", s, html)
		}
	}
}

func TestParseNorg_MissingEndFails(t *testing.T) {
	_, _, _, err := parseNorg("@document.meta\ntitle: X\n")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseNorg_UnclosedFenceFails(t *testing.T) {
	raw := `@document.meta
	title: X
	slug: x
	date: 2026-05-01
	@end
	` + "```go\nfmt.Println(1)"
	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected unclosed code fence error")
	}
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }
```

- [ ] **Step 2: Add failing loader integration assertion for checklist output**

```go
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
```

- [ ] **Step 3: Run targeted tests and confirm failure**

Run:
```bash
go test ./internal/content -run 'TestParseNorg|TestLoadPosts_SupportsNeorgTaskRendering' -v
```
Expected: FAIL with undefined `parseNorg` / missing checklist rendering.

- [ ] **Step 4: Commit failing tests**

```bash
git add internal/content/norg_parser_test.go internal/content/loader_test.go
git commit -m "test(content): add failing neorg parser coverage"
```

### Task 2: Implement Neorg parser and renderer

**Files:**
- Create: `internal/content/norg_parser.go`
- Test: `internal/content/norg_parser_test.go`

- [ ] **Step 1: Add parser entrypoint + metadata parsing**

```go
func parseNorg(raw string) (frontMatter, string, string, error) {
	meta, bodyLines, err := splitNorgFrontMatter(raw)
	if err != nil {
		return frontMatter{}, "", "", err
	}
	nodes, err := parseNorgBlocks(bodyLines)
	if err != nil {
		return frontMatter{}, "", "", err
	}
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	html := renderNorgHTML(nodes)
	return meta, body, html, nil
}
```

- [ ] **Step 2: Implement body block parser for scoped subset**

```go
type nodeKind int
const (
	nodeHeading nodeKind = iota
	nodeParagraph
	nodeUL
	nodeOL
	nodeTask
	nodeCode
)

type node struct {
	kind  nodeKind
	level int
	text  string
	items []string
	state []string
	lang  string
	code  string
}
```

Implement handlers:
- heading: `^\*+\s+`
- task: `^\*{3}\s+(TODO|DOING|DONE|CANCELLED)\s+`
- unordered list: `^-\s+`
- ordered list: `^\d+\.\s+`
- fenced code: line starts with `````; hard error on missing closing fence
- fallback paragraph for unknown non-structural lines

- [ ] **Step 3: Implement renderer with checklist attributes**

```go
func renderTaskState(s string) string {
	switch s {
	case "TODO": return "todo"
	case "DOING": return "doing"
	case "DONE": return "done"
	case "CANCELLED": return "cancelled"
	default: return "todo"
	}
}

// task block output
// <ul class="task-list"><li data-task-state="todo">...</li></ul>
```

- [ ] **Step 4: Run parser tests and confirm pass**

Run:
```bash
go test ./internal/content -run 'TestParseNorg' -v
```
Expected: PASS.

- [ ] **Step 5: Commit parser implementation**

```bash
git add internal/content/norg_parser.go internal/content/norg_parser_test.go
git commit -m "feat(content): add neorg parser and html renderer"
```

### Task 3: Integrate loader routing and keep `.md` behavior stable

**Files:**
- Modify: `internal/content/loader.go`
- Test: `internal/content/loader_test.go`

- [ ] **Step 1: Route by file extension with isolated parsers**

```go
func splitFrontMatter(raw, ext string) (frontMatter, string, string, error) {
	switch strings.ToLower(ext) {
	case ".md":
		meta, body, err := splitMarkdownFrontMatter(raw)
		return meta, body, "", err
	case ".norg":
		meta, body, html, err := parseNorg(raw)
		return meta, body, html, err
	default:
		return frontMatter{}, "", "", fmt.Errorf("unsupported content format %q", ext)
	}
}
```

- [ ] **Step 2: Use pre-rendered HTML for `.norg`, goldmark for `.md`**

```go
meta, body, norgHTML, err := splitFrontMatter(string(raw), filepath.Ext(file))
...
var htmlBody string
if strings.EqualFold(filepath.Ext(file), ".norg") {
	htmlBody = norgHTML
} else {
	htmlBody, err = renderMarkdown(body)
	if err != nil { ... }
}
```

- [ ] **Step 3: Run content test suite**

Run:
```bash
go test ./internal/content -v
```
Expected: PASS, including `.md` sort/draft tests.

- [ ] **Step 4: Commit integration**

```bash
git add internal/content/loader.go internal/content/loader_test.go
git commit -m "feat(content): route norg files to native parser"
```

### Task 4: Full verification and sample content check

**Files:**
- Optional modify: `content/posts/neorg-post.norg`

- [ ] **Step 1: Ensure sample `.norg` uses supported metadata + task syntax**

```norg
@document.meta
title: Neorg Post
slug: neorg-post
date: 2026-05-01
summary: First post using Neorg parser
tags:
  - neorg
  - notes
draft: false
@end

*** TODO polish css for task list
```

- [ ] **Step 2: Run full project tests**

Run:
```bash
go test ./...
```
Expected: PASS all packages.

- [ ] **Step 3: Manual HTML sanity check (optional)**

Run:
```bash
go test ./internal/content -run TestLoadPosts_SupportsNeorgTaskRendering -v
```
Expected: output contains `data-task-state="todo"`.

- [ ] **Step 4: Commit final polish**

```bash
git add content/posts/neorg-post.norg
git commit -m "chore(content): add sample norg task post"
```

## Self-review

- Spec coverage: all locked requirements mapped (routing, task states, supported blocks, mixed strictness, structural failure policy).
- Placeholder scan: no TBD/TODO placeholders in plan steps.
- Type consistency: uses same parser entrypoint (`parseNorg`) and same `frontMatter` model throughout.
