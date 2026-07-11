# Security Cleanup Implementation Plan

> **REQUIRED SUB-SKILL:** Use the executing-plans skill to implement this plan task-by-task.

**Goal:** Remove found security weaknesses and shorten low-risk duplication while preserving normal site behavior.

**Architecture:** Harden existing boundaries rather than adding services or dependencies. Parser validates link destinations before trusted HTML creation; HTTP layer rejects excessive search input and handles compression/statuses safely; build inputs become immutable and verified.

**Tech Stack:** Go 1.26.5, net/http, templ, GitHub Actions, Docker, Make.

---

### Task 1: Runtime and supply-chain hardening

**Files:**
- Modify: `go.mod`
- Modify: `Dockerfile`
- Modify: `Makefile`
- Modify: `.github/workflows/go-test.yml`
- Modify: `.github/workflows/release.yml`
- Modify: `README.md`

**Step 1: Update runtime references**

Replace Go 1.26.4 with 1.26.5 in Go module, Docker build image, workflows, and prerequisites documentation.

**Step 2: Verify Tailwind download**

Define SHA-256 `7d24f7fa191d2193b78cd5f5a42a6093e14409521908529f42d80b11fde1f1d4`, sourced from official v3.4.17 `sha256sums.txt`. Run `sha256sum -c` before executing downloaded binary in Docker and Make paths. Verify existing cached binary too.

**Step 3: Pin workflows and reduce permissions**

Use exact action commits with version comments:
- checkout: `34e114876b0b11c390a56381ad16ebd13914f8d5` (`v4`)
- setup-go: `7b8cf10d4e4a01d4992d18a89f4d7dc5a3e6d6f4` (`v4`)
- goreleaser action: `e435ccd777264be153ace6237001ef4d979d3a7a` (`v6`)

Add `permissions: contents: read` to test workflow. Replace GoReleaser `version: latest` with exact version `v2.17.0`. Pin Docker build and runtime images by digest.

**Step 4: Verify config**

Run: `GOTOOLCHAIN=go1.26.5 go mod verify && make test`
Expected: module verification and all tests pass.

**Step 5: Commit**

Commit message: `build: pin secure toolchain and build inputs`

### Task 2: Safe inline links

**Files:**
- Modify: `internal/content/norg_parser_test.go`
- Modify: `internal/content/norg_parser.go`

**Step 1: Write failing tests**

Add table-driven parser tests proving relative, HTTPS, HTTP, and mailto links render, while `javascript:`, `data:`, protocol-relative, malformed, and control-character destinations fail with `invalid link url`.

**Step 2: Verify red**

Run: `go test ./internal/content -run 'TestParseNorg_(ValidLinks|InvalidLinks)' -v`
Expected: invalid-link cases fail because current parser accepts them.

**Step 3: Implement minimum validation**

Add `validateLinkURL(raw string) (string, error)` using `net/url`. Trim input; reject empty, control characters, host-without-allowed-scheme, protocol-relative references, and schemes outside `http`, `https`, `mailto`. Require host for HTTP(S), address for mailto, and leading `/`, `./`, `../`, `#`, or `?` for relative links. Escape returned value only at HTML rendering.

**Step 4: Verify green**

Run focused test, then `go test ./internal/content`.
Expected: pass.

**Step 5: Commit**

Commit message: `fix: validate rendered content links`

### Task 3: HTTP resource and response hardening

**Files:**
- Modify: `internal/http/search_test.go`
- Modify: `internal/http/routes_test.go`
- Modify: `internal/http/blog_routes_test.go` if shared server test helper is needed
- Modify: `internal/http/search.go`
- Modify: `internal/http/routes.go`
- Modify: `internal/http/server.go`

**Step 1: Write failing tests**

Add tests proving:
- `/search` rejects query values over 256 bytes with `400` and does not invoke scoring.
- gzip middleware emits no body and no `Content-Encoding` for `204`.
- standard responses include `Strict-Transport-Security`, `Permissions-Policy`, and existing headers.
- constructed HTTP server has `ReadHeaderTimeout` and `MaxHeaderBytes` set.

**Step 2: Verify red**

Run: `go test ./internal/http -run 'Test(SearchRejectsOversizedQuery|CompressionSkipsNoContent|SecurityHeaders|ServerLimits)' -v`
Expected: new assertions fail against current behavior.

**Step 3: Implement minimum fixes**

Set decoded search query limit to 256 bytes and raw query limit to 512 bytes; reject oversized input with JSON `400`. Parse `Accept-Encoding` tokens and quality values, exclude bodyless statuses from compression, and close gzip writer only when compression activates. Add HSTS (`max-age=31536000; includeSubDomains`) and restrictive Permissions Policy. Configure `ReadHeaderTimeout: 5*time.Second` and `MaxHeaderBytes: 1<<20`.

**Step 4: Small safe cleanup**

Reuse one lowercase query value in `parseQuery`; preserve ranking and output.

**Step 5: Verify green**

Run focused test, `go test ./internal/http`, then `go test ./...`.
Expected: pass.

**Step 6: Commit**

Commit message: `fix: bound requests and harden responses`

### Task 4: Low-risk duplication cleanup

**Files:**
- Modify: `internal/content/loader.go`
- Modify: `internal/content/loader_extra.go`

**Step 1: Confirm baseline behavior**

Run: `go test ./internal/content -run 'TestLoad' -v`
Expected: pass.

**Step 2: Refactor**

Extract `publishedFromMeta(frontMatter) Published` and replace repeated literals for posts, projects, and project subposts. Replace identical `slugPattern`/`tagPattern` variables with one `keyPattern` used for both checks.

**Step 3: Verify behavior unchanged**

Run: `gofmt -w internal/content/loader.go internal/content/loader_extra.go && go test ./internal/content -run 'TestLoad' -v`
Expected: pass with no test changes.

**Step 4: Commit**

Commit message: `refactor: deduplicate published metadata mapping`

### Task 5: Full verification and review

**Files:**
- Modify only if verification exposes defects.

**Step 1: Generate and format**

Run: `make generate && gofmt -w cmd internal`
Expected: generation and formatting succeed.

**Step 2: Static and runtime checks**

Run:
- `go vet ./...`
- `go test ./...`
- `go test -race ./...`
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`

Expected: all pass; vulnerability scan reports no reachable vulnerabilities.

**Step 3: Diagnostics and diff review**

Run LSP diagnostics on changed Go files, inspect `git diff --check`, ensure `index.html` and `temp.html` remain untouched, and request independent security/code-quality review.

**Step 4: Fix and reverify**

Address only concrete review findings, then rerun affected checks.

**Step 5: Commit**

Commit message if needed: `test: complete security verification`
