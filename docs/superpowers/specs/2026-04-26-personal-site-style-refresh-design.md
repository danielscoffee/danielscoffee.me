# Personal Site Style Refresh Design (Soft Editorial)

**Date:** 2026-04-26
**Project:** `github.com/danielscoffee/danielscoffee.me`
**Status:** Approved for planning

---

## 1) Summary

Refresh the personal site styling with a **moderate antfu.me-inspired direction** while keeping the result simple and calm. The refresh applies to:

- Home page
- Blog index page
- Blog post detail page

Chosen direction and constraints:

- Visual direction: **B — Soft Editorial**
- Scope: **B — Home + blog list + post detail**
- Theme mode: **C — Auto light/dark + manual toggle**
- Interaction level: **B — Gentle transitions only**

No backend behavior or route structure changes are required.

---

## 2) Goals and Non-Goals

### Goals

1. Improve typography rhythm, spacing, and hierarchy for better reading comfort.
2. Introduce a subtle, polished visual identity without over-styling.
3. Add robust light/dark/system theming with a manual toggle in header.
4. Keep implementation small, maintainable, and consistent with existing Go + templ + Tailwind setup.

### Non-Goals

1. No redesign of information architecture or routing.
2. No animation-heavy interactions.
3. No client framework introduction.
4. No CMS/editor/admin feature changes.

---

## 3) Architecture and Implementation Style

Use a **hybrid styling approach**:

- Tailwind utility classes for page-level layout control.
- Small set of semantic component classes in stylesheet for consistency (header/nav/post-list/post-meta/tag chips/prose tuning/theme toggle).

This balances speed and clarity:

- Templates remain readable.
- Repeated styling concerns move into centralized CSS.
- Future tweaks stay low-risk.

---

## 4) Design Decisions

### 4.1 Visual Language

- Maintain neutral palette with slight warmth in light mode and muted coolness in dark mode.
- Increase readability through tighter typographic scale choices and improved vertical rhythm.
- Keep separators, cards, and chips low contrast and unobtrusive.

### 4.2 Layout and Components

### Header / Navigation

- Keep header simple, compact, and consistent across pages.
- Include links: Home, Blog, RSS, Theme toggle.
- Optional subtle sticky behavior is allowed if it does not increase visual noise.

### Home + Blog List

- Use breathing room between entries and cleaner metadata row.
- Add gentle link/card hover emphasis (underline/accent shift or slight background lift).
- Preserve straightforward scanning behavior.

### Post Detail

- Tune prose styles for markdown content:
  - Heading spacing and weight
  - Body line-height and measure
  - Link style consistency
  - Inline/code block legibility
  - Blockquote subtle emphasis

### Tag Chips

- Softer chip treatment than current default.
- Lower contrast backgrounds and modest rounding.
- Keep chips informative, not decorative.

### 4.3 Theme System

Theme modes:

- `system` (default)
- `light`
- `dark`

Behavior:

1. On load, resolve theme from `localStorage` (`theme-preference`) if valid.
2. If missing/invalid, use `system` and derive from `prefers-color-scheme`.
3. Manual toggle cycles modes in order: `system -> light -> dark -> system`.
4. Persist user selection to `localStorage`.
5. Update control label/aria state for accessibility.

Implementation approach:

- Theme tokens managed with CSS custom properties and/or Tailwind-compatible classes.
- Apply resolved theme early in base template to minimize flash/mismatch during initial paint.

### 4.4 Motion and Interaction

- Use short transitions (`150–200ms`) for hover/focus/theme-change affordances.
- Respect `prefers-reduced-motion` and reduce transitions accordingly.
- Avoid transforms or effects that change content structure or distract while reading.

---

## 5) File-Level Change Plan

### `internal/web/base.templ`

- Refine shell layout classes.
- Introduce theme toggle button in nav.
- Add tiny inline script for theme resolution + toggle cycle + persistence.
- Keep markup semantic and accessible.

### `internal/web/pages.templ`

- Adjust structural classes for home/blog/post pages.
- Update post list and metadata styling hooks.
- Update chip classes and prose wrapper usage.

### `internal/web/styles/input.css`

- Add base tokens for light/dark/system palette.
- Add small semantic component class group:
  - header/nav links
  - post list items
  - metadata line
  - tags/chips
  - prose overrides
  - theme toggle state
- Add reduced-motion guardrail.

### `tailwind.config.js`

- Keep scope minimal; only extend if needed for typography or color token mapping.

### Generated files

- Regenerate templ outputs and Tailwind output as part of implementation.

---

## 6) Data Flow (Theme Toggle)

1. Server renders page via templ.
2. Early client script resolves active mode (`system/light/dark`).
3. Script sets theme attribute/class on root element.
4. User presses toggle.
5. Script cycles to next mode, updates DOM state, stores preference.
6. On next page load, stored preference re-applies before interaction.

No server persistence is required.

---

## 7) Error Handling and Fallbacks

1. If `localStorage` is unavailable (privacy mode), theme logic falls back to system mode.
2. If stored value is malformed, treat as `system`.
3. If JavaScript is disabled, stylesheet still respects `prefers-color-scheme` for acceptable baseline.
4. Theme toggle should fail safe (site remains readable and navigable).

---

## 8) Testing Strategy

### 8.1 Automated

- `templ generate -path .`
- `go test ./...`
- `go build ./...`

### 8.2 Manual UI Checks

1. Verify visual hierarchy on home/blog/post pages in both light and dark.
2. Verify toggle cycle order and persistence across reloads.
3. Verify keyboard accessibility for toggle and nav links.
4. Verify no harsh motion when `prefers-reduced-motion` is enabled.
5. Verify markdown readability for headings, paragraphs, code, blockquotes, and links.
6. Verify mobile and desktop layouts remain clean and uncluttered.

---

## 9) Risks and Mitigations

1. **Risk:** Theme flash on load.
   - **Mitigation:** Apply resolved mode in early inline script before main paint-sensitive styling.

2. **Risk:** Over-styling drifts from “simple” requirement.
   - **Mitigation:** Keep palette neutral and motion restrained; avoid decorative effects.

3. **Risk:** Template class clutter.
   - **Mitigation:** Move repeated patterns into small semantic classes.

---

## 10) Acceptance Criteria

1. Site clearly reflects **Soft Editorial** style without visual excess.
2. Home, blog list, and post detail pages share consistent spacing/type language.
3. `system/light/dark` toggle works, persists, and remains accessible.
4. Interactions are gentle and reduced-motion aware.
5. Build/test pipeline remains green.

---

## 11) Out of Scope for This Pass

- Search redesign
- Comments/analytics UI updates
- Content model changes
- New page types
- Complex motion system
