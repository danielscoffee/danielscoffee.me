# Security Cleanup Design

## Scope

Apply targeted security hardening plus safe code simplification without changing site routes, content format, or normal rendering. Preserve generated-file workflow and ignore untracked `index.html`/`temp.html`.

## Security changes

1. Upgrade Go 1.26.4 to 1.26.5 everywhere. `govulncheck` reports reachable `GO-2026-5856` in Go 1.26.4 `crypto/tls`; Go 1.26.5 fixes it.
2. Validate rendered inline-link destinations before placing parser output into trusted HTML. Allow relative URLs and `http`, `https`, and `mailto`; reject protocol-relative URLs and dangerous/unknown schemes.
3. Bound search query input before tokenization/scoring to prevent request-cost amplification. Oversized HTTP queries return `400`.
4. Fix gzip middleware so uncompressed statuses do not flush an empty gzip stream. Add HSTS and Permissions Policy response headers plus explicit HTTP header timeout/size limits.
5. Pin downloaded Tailwind binary by official SHA-256, pin GitHub Actions to commit SHAs, reduce test workflow token permissions, and pin GoReleaser major version.

## Simplification

Extract repeated `frontMatter` to `Published` field copying into one helper and reuse identical slug/tag validation regex. Avoid parser rewrites and clever abstractions: risk exceeds line savings.

## Validation

Use TDD for Go behavior changes: add focused failing tests, observe expected failures, implement minimum fixes, then rerun package and full tests. Run `gofmt`, `go vet`, `go test -race ./...`, `govulncheck ./...`, LSP diagnostics, and final independent review. Config-only changes receive syntax/build verification rather than artificial unit tests.
