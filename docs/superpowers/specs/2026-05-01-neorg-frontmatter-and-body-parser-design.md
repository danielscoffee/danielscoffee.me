# Neorg Frontmatter + Body Parser Design

Date: 2026-05-01
Status: Draft for review
Owner: content pipeline

## Goal

Add real `.norg` support with two content modes only:

- `.md` uses Markdown + YAML frontmatter (`---` blocks)
- `.norg` uses Neorg metadata (`@document.meta` ... `@end`) and Neorg body parser

Neorg body must parse task syntax like `*** TODO`, not treat it as plain Markdown.

## Locked Decisions

1. Strictness: mixed
   - strict for frontmatter + core structure
   - lenient fallback for unsupported non-breaking body syntax
2. Task rendering: checklist HTML with state attributes
3. Task states v1: `TODO`, `DOING`, `DONE`, `CANCELLED`
4. Body scope v1: headings, paragraphs, unordered/ordered lists, fenced code, links
5. Parser source: in-house parser
6. Unsupported constructs policy:
   - hard error only on structural break
   - otherwise render escaped paragraph fallback

## Architecture

### High-level

- Keep `internal/content/loader.go` as extension router.
- Add `internal/content/norg_parser.go` for `.norg` parsing/rendering.
- Keep `.md` path unchanged.

### Routing

- `.md` -> existing Markdown parser path
- `.norg` -> Neorg parser path
- any other extension -> unsupported format error

### Neorg parse pipeline

1. Split + parse frontmatter (`@document.meta` ... `@end`) strict
2. Parse body lines into small block AST
3. Render AST into HTML
4. Return body raw text + rendered HTML

## Components

## 1) `splitNorgFrontMatter(raw string) (frontMatter, []string, error)`

Responsibilities:
- require first non-empty line to be `@document.meta`
- find matching `@end`
- parse key/value metadata lines
- validate required fields: `title`, `slug`, `date`
- return remaining body lines

Metadata parsing rules:
- accept `key: value`
- `tags` supports:
  - `tags: a, b`
  - `tags: [a, b]`
  - multiline list:
    - `tags:` then `- item` lines
- `draft` parse bool (`true/false`)

Strict failures:
- missing `@document.meta` / missing `@end`
- malformed required metadata lines
- invalid `draft`

## 2) `parseNorgBlocks(lines []string) ([]Node, error)`

Node types:
- `Heading{Level int, Text string}`
- `Paragraph{Text string}`
- `List{Ordered bool, Items []ListItem}`
- `TaskList{Items []TaskItem}`
- `CodeBlock{Lang string, Code string}`

Task item:
- `TaskItem{State string, Text string}`
- states allowed: TODO, DOING, DONE, CANCELLED

Block syntax v1:
- Heading: one or more `*` then space then text
- Task: `*** <STATE> <text>`
- UL: `- text`
- OL: `1. text`
- Code fence: triple backticks with optional lang
- Paragraph: contiguous plain lines

Lenient fallback:
- unsupported single lines => escaped paragraph

Hard errors (structural):
- unclosed code fence
- impossible nesting transition requiring recovery that changes block boundaries

## 3) `renderNorgHTML(nodes []Node) string`

Render rules:
- headings -> `<h1..h6>`
- paragraphs -> `<p>` escaped text + inline link render
- UL/OL -> `<ul>/<ol><li>`
- task list -> `<ul class="task-list">`
  - each item: `<li data-task-state="todo|doing|done|cancelled">...`
- code fences -> `<pre><code class="language-...">...`

Inline handling v1:
- links `[text](url)` recognized
- all other inline markup escaped/pass-through

## Data Flow

1. `LoadPosts` reads file
2. extension route:
   - `.md`: existing path
   - `.norg`: new Neorg path
3. parse frontmatter
4. parse + render body
5. append `Post` with:
   - `BodyHTML`: rendered output
   - `BodyMD`: raw body text (legacy field name retained)

## Error Handling Contract

Return `parse <file>: <reason>` from loader.

Reasons include:
- `missing frontmatter delimiter`
- `invalid frontmatter structure`
- `invalid neorg metadata line "..."`
- `invalid draft value "..."`
- `title, slug, and date are required`
- `invalid task state "..."`
- `unclosed code fence`

## Testing Strategy

Add `internal/content/norg_parser_test.go`:

1. frontmatter parse success
2. missing `@end` -> error
3. task states render to `data-task-state`
4. headings/lists/code/links happy path
5. unsupported line fallback paragraph
6. unclosed fence -> error
7. invalid task state -> error

Update `internal/content/loader_test.go`:
- verify `.md` still works unchanged
- verify `.norg` route uses Neorg parser behavior for tasks

## Trade-offs

Chosen: single-pass line parser.

Pros:
- small code surface
- easy debug
- low dependency risk

Cons:
- limited inline feature set in v1
- deeper nesting/extensions later need careful parser growth

Not chosen:
- two-pass tokenizer/parser (cleaner long-term, more initial complexity)
- regex-only transform (faster initially, brittle semantics)

## Out of Scope (v1)

- full Neorg spec coverage
- rich inline formatting beyond links
- embedded directives/macros outside chosen subset
- task dependency/priority metadata

## Acceptance Criteria

- `.md` path unchanged and passing tests
- `.norg` supports defined metadata and body subset
- `*** TODO/DOING/DONE/CANCELLED` render checklist HTML with state attributes
- structural syntax errors fail fast
- unsupported non-structural syntax degrades to paragraph fallback
- tests cover parser + loader integration

## Implementation Notes

- keep parser package-local under `internal/content`
- keep escape handling centralized to avoid HTML injection
- avoid silent drops of text
- keep error messages deterministic for tests
