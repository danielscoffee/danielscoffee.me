# Norg Parser Package Implementation Plan

> **REQUIRED SUB-SKILL:** Use the executing-plans skill to implement this plan task-by-task.

**Goal:** Move Norg parsing into focused `internal/content/norg` files while preserving parser and loader behavior.

**Architecture:** `norg.Parse` becomes package entrypoint and returns validated `norg.Meta`, original body, and safe rendered HTML. Existing content loaders consume that result; parser internals remain package-private and split by responsibility.

**Tech Stack:** Go 1.26.5, standard library, Chroma v2, existing Go tests.

---

### Task 1: Establish new package contract with characterization tests

**Files:**
- Move: `internal/content/norg_parser_test.go` → `internal/content/norg/parser_test.go`

**Step 1: Move test file**

Run:

```bash
mkdir -p internal/content/norg
git mv internal/content/norg_parser_test.go internal/content/norg/parser_test.go
```

Change package declaration to `package norg`. Rename calls from `parseNorg` to desired package entrypoint `Parse`; keep private-helper tests unchanged because tests remain in same package.

**Step 2: Run test to verify RED**

Run:

```bash
go test ./internal/content/norg
```

Expected: build failure for undefined `Parse`, `parseTagValues`, `renderNorgHTML`, and parser node types. Failure proves tests now target new package.

### Task 2: Split parser implementation into package files

**Files:**
- Create: `internal/content/norg/parser.go`
- Create: `internal/content/norg/metadata.go`
- Create: `internal/content/norg/blocks.go`
- Create: `internal/content/norg/inline.go`
- Create: `internal/content/norg/render.go`

**Step 1: Add parser entrypoint and node model**

`parser.go` defines private node kinds/types and:

```go
func Parse(raw string) (Meta, string, string, error) {
	meta, bodyLines, err := splitFrontMatter(raw)
	if err != nil {
		return Meta{}, "", "", err
	}
	nodes, err := parseBlocks(bodyLines)
	if err != nil {
		return Meta{}, "", "", err
	}
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	rendered, err := renderHTML(nodes)
	if err != nil {
		return Meta{}, "", "", err
	}
	return meta, body, rendered, nil
}
```

**Step 2: Move metadata parsing and validation**

`metadata.go` defines exported `Meta`, private `splitFrontMatter`, metadata/tag parsing, quote stripping, and existing slug/date/tag validation. Preserve error text exactly.

**Step 3: Move block parsing**

`blocks.go` contains `parseBlocks`, block collection, structural line parsers, task states, and image alt extraction. Rename only top-level helpers needed for package clarity; avoid generalized parser abstractions.

**Step 4: Move inline parsing and URL validation**

`inline.go` contains modifier rendering, links/images, carryover attributes, and URL validators. Preserve control-character, scheme, CDN-host, and path checks.

**Step 5: Move HTML/code/table rendering**

`render.go` contains node rendering, Chroma fallback, plain code output, and table parsing/rendering. Replace task-state switch with `strings.ToLower(item.state)` only because block parser already validates states.

**Step 6: Format and verify GREEN**

Run:

```bash
gofmt -w internal/content/norg/*.go
go test ./internal/content/norg
```

Expected: all 14 parser tests pass.

**Step 7: Commit package extraction**

```bash
git add internal/content/norg
git commit -m "refactor: extract norg parser package"
```

### Task 3: Wire content loaders to package

**Files:**
- Modify: `internal/content/loader.go`
- Modify: `internal/content/loader_extra.go`
- Delete: `internal/content/norg_parser.go`
- Test: `internal/content/loader_test.go`

**Step 1: Import package and use metadata type**

In `loader.go`, import:

```go
"github.com/danielscoffee/danielscoffee.me/internal/content/norg"
```

Change `publishedFromMeta`, `splitFrontMatter`, and `contentEntry.meta` to use `norg.Meta`. Delegate `.norg` parsing to `norg.Parse` and preserve unsupported-extension error.

**Step 2: Remove duplicate validation**

Delete `frontMatter`, `validateFrontMatter`, and its regex/time imports from `loader.go`. Delete two post-parse `validateFrontMatter` calls from `loader_extra.go`; `norg.Parse` already validates every loaded file.

**Step 3: Delete old implementation**

```bash
git rm internal/content/norg_parser.go
```

**Step 4: Run focused integration tests**

Run:

```bash
gofmt -w internal/content/loader.go internal/content/loader_extra.go
go test ./internal/content/... -v
```

Expected: parser and loader tests pass, including invalid metadata, projects, pages, tasks, and rendered HTML.

**Step 5: Commit loader wiring**

```bash
git add -A internal/content
git commit -m "refactor: use norg parser package"
```

### Task 4: Verify repository and reduction

**Files:**
- Verify all modified files

**Step 1: Run proactive diagnostics**

Run LSP diagnostics on `internal/content` and resolve new errors only.

**Step 2: Run focused checks**

```bash
go test ./internal/content/...
go vet ./internal/content/...
```

Expected: pass.

**Step 3: Run generated full suite**

```bash
make test
go test -race ./...
go vet ./...
```

Expected: pass. `make test` regenerates ignored templ files required by full build.

**Step 4: Inspect scope and line count**

```bash
git diff --check
git status --short
wc -l internal/content/norg/*.go internal/content/loader.go internal/content/loader_extra.go
```

Expected: no whitespace errors; only planned files and plan docs changed. Parser implementation split across five focused files, with total code reduced where duplication was removed.

**Step 5: Review final diff**

Confirm exported surface is only `Meta` and `Parse`; errors, HTML, validation, accessibility attributes, and Chroma fallback remain unchanged. Request read-only review, route any fixes through original writer, then rerun affected checks.

**Step 6: Commit documentation or final fixes**

```bash
git add docs/plans/2026-07-25-norg-parser-package-design.md docs/plans/2026-07-25-norg-parser-package.md
git commit -m "docs: record norg parser refactor plan"
```
