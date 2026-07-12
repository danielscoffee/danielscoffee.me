# Editorial Visual Refresh Implementation Plan

> **REQUIRED SUB-SKILL:** Use the executing-plans skill to implement this plan task-by-task.

**Goal:** Deliver warm, accessible editorial polish across every rendered page while preserving application behavior.

**Architecture:** Build visual system from semantic CSS tokens and reusable templ components. Change shell and page markup only where current structure blocks hierarchy or accessibility; keep search/theme JavaScript small and dependency-free.

**Tech Stack:** Go, templ, Tailwind CSS 3.4.17, vanilla JavaScript, self-hosted WOFF2 fonts.

---

### Task 1: Capture baseline and add font assets

**Files:**
- Create: `internal/web/assets/fonts/source-serif-4-regular.woff2`
- Create: `internal/web/assets/fonts/source-serif-4-semibold.woff2`
- Create: `internal/web/assets/fonts/source-sans-3-regular.woff2`
- Create: `internal/web/assets/fonts/source-sans-3-semibold.woff2`
- Create: `internal/web/assets/fonts/OFL.txt`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/web/styles_test.go`

**Step 1: Capture baseline**

Run application and capture light/dark screenshots for `/blog`, `/post/<slug>`, `/projects`, `/projects/<slug>`, `/about`, and search at 375px and 1440px. Store outside repository as review evidence.

**Step 2: Write failing asset/style tests**

Assert font files and license exist. Assert CSS contains named font faces and semantic typography tokens.

Run: `go test ./internal/web -run 'Test(FontAssets|InputStyles)' -v`
Expected: FAIL because assets/tokens do not exist.

**Step 3: Add pinned self-hosted fonts**

Add Source Serif 4 and Source Sans 3 WOFF2 files from official OFL-licensed release. Record source release and SHA-256 values in commit message or plan notes.

**Step 4: Add font faces and tokens**

Define regular/semibold `@font-face` declarations with `font-display: swap`. Add `--font-editorial`, `--font-ui`, and `--font-code` tokens without changing components yet.

**Step 5: Verify and commit**

Run: `make generate && go test ./internal/web -v`
Expected: PASS.

Commit: `style: add self-hosted editorial typography`

### Task 2: Establish warm visual tokens and base behavior

**Files:**
- Modify: `internal/web/styles/input.css`
- Modify: `internal/web/styles_test.go`

**Step 1: Write failing token tests**

Assert light/dark themes define background, surface, elevated surface, text, muted, border, accent, focus, shadow, radius, and width tokens. Assert global `:focus-visible`, selection, and reduced-motion hooks.

**Step 2: Verify red**

Run: `go test ./internal/web -run TestInputStyles -v`
Expected: FAIL on missing tokens/hooks.

**Step 3: Implement tokens and base styles**

Replace gray palette with cream/espresso/caramel/rust values. Use roasted-brown dark mode. Apply UI font to shell, editorial font to long-form content, consistent focus ring, selection color, and motion defaults.

**Step 4: Verify and commit**

Run: `make generate && go test ./internal/web -v`
Expected: PASS.

Commit: `style: establish warm editorial design tokens`

### Task 3: Refine site shell and navigation

**Files:**
- Modify: `internal/web/base.templ`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/http/blog_routes_test.go`

**Step 1: Write failing rendered-markup tests**

Assert pages contain skip link targeting `#main-content`, main landmark ID, masthead structure, accessible theme/search labels, dialog close button, and footer with RSS link.

**Step 2: Verify red**

Run: `go test ./internal/http -run 'Test(SiteShell|SearchDialog)' -v`
Expected: FAIL on missing landmarks/controls.

**Step 3: Implement shell markup**

Reorganize header into brand, primary nav, and compact actions. Add skip link, main ID, search close button, and restrained footer. Keep all current destinations and script loading.

**Step 4: Implement responsive shell CSS**

Use contained masthead, predictable mobile wrapping, minimum 44px controls, visible focus, and purpose-specific shell widths. Avoid horizontal scrolling in action row.

**Step 5: Verify and commit**

Run: `make generate && go test ./internal/http ./internal/web -v`
Expected: PASS.

Commit: `style: refine responsive site shell`

### Task 4: Consolidate page intros and index cards

**Files:**
- Modify: `internal/web/pages.templ`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/http/blog_routes_test.go`
- Modify: `internal/http/project_routes_test.go`

**Step 1: Write failing markup tests**

Assert blog/projects/tag pages render reusable intro classes, card-list structure, semantic `time datetime`, and existing links/tags.

**Step 2: Verify red**

Run: `go test ./internal/http -run 'Test(BlogRoutes|ProjectRoutes|EditorialIndexes)' -v`
Expected: FAIL only on new structure/semantics.

**Step 3: Extract templ components**

Create `PageIntro`, shared metadata, and tag-list components. Keep page titles/descriptions/content unchanged. Consolidate post/project card markup where types permit without awkward abstractions.

**Step 4: Style index cards**

Use wider index container, consistent card spacing, readable summary measure, stable tag/date grouping, and restrained hover/focus elevation.

**Step 5: Verify and commit**

Run: `templ generate -path . && go test ./internal/http ./internal/web -v`
Expected: PASS.

Commit: `style: unify editorial index cards`

### Task 5: Improve article reading experience

**Files:**
- Modify: `internal/web/pages.templ`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/http/blog_routes_test.go`
- Modify: `internal/http/project_routes_test.go`

**Step 1: Write failing structure tests**

Assert article pages separate header and body, render semantic dates, and preserve raw parsed content inside article body container.

**Step 2: Verify red**

Run focused post/project route tests.
Expected: FAIL on new containers/semantics.

**Step 3: Refine templates**

Reuse article header component for posts, projects, subposts, and about where appropriate. Preserve title, date, summary, breadcrumb, and content.

**Step 4: Refine prose CSS**

Set body measure near 68ch. Improve heading rhythm, paragraph spacing, links, lists, blockquotes, definitions, task states, tables, code, and images. Add mobile overflow containment for tables/pre blocks.

**Step 5: Verify and commit**

Run: `make generate && go test ./...`
Expected: PASS.

Commit: `style: improve long-form reading experience`

### Task 6: Simplify project devlog grid

**Files:**
- Modify: `internal/web/pages.templ`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/web/styles_test.go`
- Modify: `internal/http/project_routes_test.go`

**Step 1: Write failing tests**

Assert project subposts use semantic stable card grid and CSS no longer contains count-specific selectors such as `[data-count="3"]` or `nth-child(6n+1)`.

**Step 2: Verify red**

Run project/style tests.
Expected: FAIL on current count-dependent grid.

**Step 3: Simplify markup and CSS**

Remove `strconv` and `data-count`. Use one-column mobile and two-column desktop grid with consistent cards, dates, summaries, focus, and hover states.

**Step 4: Verify and commit**

Run: `make generate && go test ./internal/http ./internal/web -v`
Expected: PASS.

Commit: `refactor: simplify project devlog cards`

### Task 7: Polish search and theme controls

**Files:**
- Modify: `internal/web/base.templ`
- Modify: `internal/web/assets/js/search.js`
- Modify: `internal/web/assets/js/theme-toggle.js`
- Modify: `internal/web/styles/input.css`
- Modify: `internal/web/styles_test.go`

**Step 1: Write failing source/markup tests**

Assert search keeps `innerHTML` forbidden, uses `AbortController`, exposes loading/error state, renders messages as non-links, and close button closes dialog. Assert theme button retains full accessible state while visible text stays compact.

**Step 2: Verify red**

Run: `go test ./internal/web -run 'Test(Search|Theme)' -v`
Expected: FAIL on new behavior markers.

**Step 3: Implement minimal JS changes**

Cancel stale search requests, render loading/empty/error states as text nodes, close with explicit button, preserve Ctrl/Cmd+K and Escape, and restore focus to trigger. Keep theme storage/mode behavior unchanged.

**Step 4: Style controls/modal**

Apply command-palette surface, clear result hierarchy, mobile-safe dimensions, loading/error styling, and compact fixed-width theme/search controls.

**Step 5: Verify and commit**

Run: `make generate && go test ./internal/web ./internal/http -v`
Expected: PASS.

Commit: `style: polish search and theme controls`

### Task 8: Visual, accessibility, and regression pass

**Files:**
- Modify only files where review finds concrete defects.
- Generated: `internal/web/assets/css/output.css`

**Step 1: Generate and verify code**

Run:
- `make generate`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- LSP diagnostics on changed Go/templ files
- `git diff --check`

Expected: all pass.

**Step 2: Run visual matrix**

Review blog, post, projects, project detail, devlog, about, search, 404, and 500 at 375px, 768px, and 1440px in light/dark/system themes.

**Step 3: Run accessibility matrix**

Check keyboard-only navigation, focus return, dialog behavior, 200% zoom, reduced motion, contrast, long text/tags, table/code overflow, and screen-reader landmark order.

**Step 4: Compare scope and size**

Confirm no backend/parser/API changes. Confirm obsolete CSS removed and source remains organized into one components layer. Document remaining visual risks.

**Step 5: Final review and commit**

Request independent spec and code-quality review; fix concrete findings and rerun affected checks.

Commit if needed: `style: complete editorial visual refresh`
