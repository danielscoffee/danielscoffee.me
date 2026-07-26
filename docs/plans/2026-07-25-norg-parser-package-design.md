# Norg Parser Package Design

## Goal

Replace 1,015-line `internal/content/norg_parser.go` with focused files under existing empty `internal/content/norg/` package. Preserve content format, rendered HTML, errors, validation, and loader APIs.

## Package boundary

`internal/content/norg` owns Norg metadata parsing, validation, block parsing, inline rendering, URL safety, tables, and syntax highlighting. Package exports only:

```go
type Meta struct {
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Draft   bool
}

func Parse(raw string) (Meta, string, string, error)
```

All parser implementation remains private. `internal/content` imports package and converts `norg.Meta` into existing domain models. No HTTP or web package changes.

## Files

- `parser.go`: `Parse` orchestration and parser node types.
- `metadata.go`: frontmatter extraction, metadata parsing, tag parsing, and validation.
- `blocks.go`: block scanner plus heading, list, task, quote, definition, code, table, and image-line recognition.
- `inline.go`: inline modifiers, links, images, carryover attributes, and URL validation.
- `render.go`: node HTML rendering, tables, task state output, and Chroma code highlighting.

Tests move beside package implementation. Existing loader tests remain in `internal/content` and cover integration.

## Reduction

Remove old parser file, private duplicate metadata type, and redundant project metadata validation after parsing already succeeded. Collapse trivial duplicate helpers only where behavior stays obvious. Keep parser AST, Chroma dependency, security checks, accessibility attributes, and exact errors. No generalized parser framework, interfaces, registries, or new dependencies.

## Data flow

`loadContentFile` reads `.norg` file, `splitFrontMatter` delegates to `norg.Parse`, and returned `norg.Meta`, source body, and trusted rendered HTML populate existing content models. Parser validates metadata before returning. Inline renderer validates links and CDN image URLs before emitting HTML.

## Validation

Move existing characterization tests before implementation and observe new package test failure. Make package tests pass, then wire loaders and remove old files. Run `gofmt`, focused package tests, full generated test suite, race tests, vet, diagnostics, and diff checks. Compare parser line count before and after without accepting behavior loss solely for fewer lines.
