# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Five audiences, all confirmed, none dominant:

- **Backend engineers and peers** reading posts and devlogs on Go, distributed systems, infrastructure, devtools, and event-driven systems.
- **Hiring managers and recruiters** evaluating Daniel as a backend engineer through his writing and shipped projects.
- **Potential collaborators** (open source, side projects) who arrive through a project page or GitHub.
- **Clients** looking for someone to build or fix backend systems.
- **Daniel himself**: the site is also a public notebook and learning log. Writing for future-self is a legitimate use, not a fallback.

Most visitors arrive cold on a single post or project from a link, search, or RSS, then decide whether to look around.

## Product Purpose

Personal site and custom blog engine for Daniel (danielscoffee), a backend developer. It publishes blog posts, project pages, and nested project devlogs, with search, tags, RSS, and a sitemap.

Success, all confirmed:

- **Proof of work.** Projects and devlogs show real, shipped engineering that backs up every claim.
- **Steady readers.** People come back or subscribe by RSS to follow the writing.
- **The engine itself.** The site is a showcase of its own stack: Go, templ, and a hand-written Neorg (`.norg`) parser.
- **Get contacted.** Visitors reach out through GitHub or LinkedIn, including client and collaboration enquiries.

## Positioning

Writing and devlogs from someone who works on production healthcare systems (Java, PL/SQL, Kafka, AWS, Docker, event-driven systems) and builds Go infrastructure, routing engines, devtools, and distributed-systems experiments after hours. The engine is self-built rather than a hosted platform, and posts are authored in Neorg.

## Operating Context

- Content is authored as `.norg` files in `content/` (posts, pages, projects with devlog subposts) and loaded into memory at startup.
- Visitors read long-form technical prose with code blocks, tables, task lists, and definition lists.
- Navigation runs through Blog, Projects, and About, plus Ctrl+K search, tag pages, and RSS.
- Light, dark, and system themes are persisted in `localStorage`.
- Deployed via Docker on Railway; releases cut through GoReleaser.

## Capabilities and Constraints

- Server-rendered Go `net/http` + templ + Tailwind CSS v3 (standalone binary); assets embedded via `embed.FS`.
- `output.css` and `*_templ.go` are generated; edit `input.css` and `.templ` sources.
- Strict CSP: `default-src 'self'`, `script-src 'self'` (no inline scripts), `img-src 'self' https: data:`, `font-src 'self'`.
- No external runtime font or image requests; fonts are self-hosted WOFF2.
- JavaScript is minimal (theme init/toggle, search); DOM rendering never uses `innerHTML`.
- Content images only from HTTPS hosts containing `cdn`.
- Routes, content schema, and the search API are stable contracts.
- **Open:** no contact form, email address, service list, rates, or availability statement exists. Client-facing offers are undecided and must not be invented.

## Brand Commitments

- **Name:** danielscoffee.
- **Coffee is core identity**, not a pun: coffee is a real part of Daniel's life and may recur as a motif across the site.
- **Voice:** dry, terse, self-deprecating backend humor ("Backend developer. Systems enjoyer. Production bug survivor."). Values software that is "simple, observable, fast enough, and boring in production." Closes with "Still early. Still learning. Still shipping."

## Evidence on Hand

- About page: `content/pages/about.norg`.
- Posts: `content/posts/` (`blogcreation.norg` is nearly empty).
- Project: `content/projects/blog-engine/` with devlogs (`how-i-rebuilt`, `search-ranking`, `devlog-20260519`).
- Profiles: GitHub `github.com/danielscoffee`, LinkedIn `linkedin.com/in/daniel-reis-santos`.
- **Absent, do not fabricate:** testimonials, client names, employer names, metrics, case studies, photos of Daniel, press.

## Product Principles

1. **Shipped over claimed.** Every statement about skill points at a post, project, or devlog that proves it.
2. **Boring in production.** The site practices what the writing preaches: fast, server-rendered, minimal JS, no external runtime dependencies.
3. **Writing is the product.** Reading comfort for long technical prose outranks decoration on content pages.
4. **One identity, many visitors.** Peers, recruiters, collaborators, and clients all meet the same honest Daniel; no audience-specific personas.
5. **Coffee is who, not what.** The coffee identity frames the site; it never stands in for the engineering content.

## Accessibility & Inclusion

Established in the 2026-07-12 editorial refresh and binding: skip link, visible focus, 44px minimum touch targets, 200% zoom support, reduced-motion respected, no color-only state, keyboard-operable search dialog.
