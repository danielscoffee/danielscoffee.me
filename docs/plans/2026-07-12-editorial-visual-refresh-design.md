# Editorial Visual Refresh Design

## Goal

Polish existing site into warm coffee/editorial experience while preserving routes, content format, search API, theme modes, and server behavior. Improve hierarchy, readability, responsive behavior, accessibility, and visual consistency without full redesign.

## Direction

Use warm paper backgrounds, espresso text, caramel surfaces, and restrained rust accents. Dark mode uses roasted-brown surfaces rather than neutral black. Self-host Source Serif 4 for editorial headings/prose and Source Sans 3 for navigation, metadata, and controls; keep system monospace for code.

Use purpose-specific widths: narrow article body around 68 characters, medium content shell for article furniture, and wider index shell for post/project cards. Replace flat lists and count-dependent project geometry with predictable responsive editorial cards. Keep decoration subtle: borders, small shadows, paper-like tonal layers, and no ornamental imagery.

## Components

- Header: clearer masthead, primary navigation, compact search/theme/social actions, visible focus states.
- Footer: restrained RSS/social links and site identity.
- Index intros: reusable eyebrow/title/subtitle hierarchy.
- Post/project cards: consistent title, summary, date, and tags.
- Articles: narrow readable measure, stronger heading rhythm, refined links, blockquotes, tables, code, tasks, definitions, and images.
- Project devlogs: stable one/two-column responsive grid without item-count selectors.
- Search: command-palette presentation with loading, empty, error, and results states.

## Accessibility and performance

Add skip link, semantic `time` elements, dialog close control, minimum touch targets, global visible focus treatment, non-color-only states, 200% zoom support, and reduced-motion preservation. Self-host only required WOFF2 weights and preload critical fonts. No external runtime font or image requests.

## Scope boundaries

No backend, route, content schema, or search API changes. Accessibility exception: only parser change offsets Norg headings by one level (capped at `h6`) to reserve `h1` for page title. Generated `output.css` remains generated-only. JavaScript changes stay limited to search state/focus behavior and compact theme labels. Preserve safe DOM rendering without `innerHTML`.

## Validation

Review blog, article, projects, project detail, devlog, about, search, and error pages at 375px, 768px, and 1440px in light/dark/system themes. Check keyboard flow, focus visibility, overflow, long titles/tags, code, tables, 200% zoom, and reduced motion. Run generation, Go tests, race tests, vet, diagnostics, and independent review.
