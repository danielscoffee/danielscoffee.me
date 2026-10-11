---
name: danielscoffee
description: A backend engineer's warm-paper notebook, brewed in Go, templ, and Neorg.
colors:
  # Light theme (default, :root)
  oat-paper: "rgb(250 246 238)"
  filter-paper: "rgb(255 251 244)"
  fresh-filter: "rgb(255 253 248)"
  espresso-ink: "rgb(54 35 25)"
  cold-brew: "rgb(105 82 68)"
  kraft-line: "rgb(213 193 174)"
  cinnamon-roast: "rgb(159 75 35)"
  cinnamon-roast-deep: "rgb(121 52 24)"
  ember-focus: "rgb(173 81 32)"
  latte-foam: "rgb(241 229 211)"
  mocha-label: "rgb(91 58 40)"
  parchment-code: "rgb(239 229 215)"
  code-ink: "rgb(66 42 29)"
  roast-shadow: "rgb(63 40 27)"
  kraft-frame: "rgb(176 132 88)"
  # Dark theme (:root[data-theme="dark"] and system-dark)
  roasted-ground: "rgb(27 20 16)"
  dark-roast: "rgb(38 28 23)"
  medium-roast: "rgb(50 37 29)"
  cream: "rgb(244 231 211)"
  milky-tan: "rgb(190 166 143)"
  burnt-kraft: "rgb(91 70 57)"
  caramel-crema: "rgb(226 151 86)"
  caramel-crema-light: "rgb(244 184 119)"
  caramel-focus: "rgb(236 159 91)"
  dark-latte-foam: "rgb(61 44 34)"
  dark-mocha-label: "rgb(225 199 171)"
  grounds-code: "rgb(22 16 13)"
  dark-code-ink: "rgb(238 216 190)"
typography:
  brand:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "clamp(1.25rem, 1.1rem + 0.5vw, 1.45rem)"
    fontWeight: 600
  crumb:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "0.95rem"
    fontWeight: 400
  headline:
    fontFamily: "Source Sans 3, ui-sans-serif, system-ui, sans-serif"
    fontSize: "2.25rem"
    fontWeight: 600
    lineHeight: 1.2
    letterSpacing: "-0.02em"
  title:
    fontFamily: "Source Sans 3, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: "-0.02em"
  board-title:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "clamp(2rem, 5vw, 2.5rem)"
    fontWeight: 600
    lineHeight: 1.1
    letterSpacing: "-0.015em"
  board-row:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "clamp(1.5rem, 4vw, 2rem)"
    fontWeight: 600
    lineHeight: 1.15
    letterSpacing: "-0.015em"
  card-title:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "clamp(1.3rem, 3vw, 1.55rem)"
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: "-0.015em"
  body-prose:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "clamp(1.03rem, 0.99rem + 0.2vw, 1.12rem)"
    fontWeight: 400
    lineHeight: 1.82
  body:
    fontFamily: "Source Sans 3, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.6
    fontFeature: "\"rlig\" 1, \"calt\" 1"
  label:
    fontFamily: "Source Sans 3, ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.82rem"
    fontWeight: 500
  mono:
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace"
    fontSize: "0.875rem"
    lineHeight: 1.5
rounded:
  sm: "0.375rem"
  md: "0.75rem"
spacing:
  gutter: "1rem"
  grid-gap: "1.25rem"
  section-gap: "2.5rem"
  touch-target: "44px"
  content-reading: "68ch"
  content-wide: "72rem"
  plate-rest: "96px"
  plate-touch: "128px"
  plate-open: "192px"
components:
  brand:
    textColor: "{colors.espresso-ink}"
    typography: "{typography.brand}"
    height: "44px"
  brand-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  nav-link:
    textColor: "{colors.cold-brew}"
    typography: "{typography.label}"
    height: "44px"
    padding: "0.35rem 0.55rem"
  nav-link-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  button-ghost:
    textColor: "{colors.cold-brew}"
    typography: "{typography.label}"
    height: "44px"
    padding: "0.35rem 0.55rem"
  button-ghost-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  button-outline:
    textColor: "{colors.espresso-ink}"
    rounded: "{rounded.sm}"
    size: "44px"
  index-entry:
    textColor: "{colors.espresso-ink}"
  index-entry-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  index-entry-meta:
    textColor: "{colors.cold-brew}"
  article-tag:
    textColor: "{colors.cold-brew}"
  article-tag-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  input-search:
    backgroundColor: "{colors.oat-paper}"
    textColor: "{colors.espresso-ink}"
    rounded: "{rounded.sm}"
    height: "48px"
    padding: "0.7rem 0.85rem"
  menu-board:
    textColor: "{colors.espresso-ink}"
    width: "100%"
  board-title:
    textColor: "{colors.cinnamon-roast-deep}"
    typography: "{typography.board-title}"
  board-row:
    textColor: "{colors.espresso-ink}"
    typography: "{typography.board-row}"
    height: "44px"
  scene-plate:
    backgroundColor: "{colors.kraft-frame}"
    size: "96px"
  crumb-last:
    textColor: "{colors.cinnamon-roast-deep}"
    typography: "{typography.crumb}"
    height: "44px"
  site-section-hover:
    textColor: "{colors.espresso-ink}"
  crumb-link:
    textColor: "{colors.cold-brew}"
    height: "44px"
  crumb-link-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  dialog-search:
    backgroundColor: "{colors.fresh-filter}"
    textColor: "{colors.espresso-ink}"
    rounded: "{rounded.md}"
    padding: "1.15rem"
---

# Design System: danielscoffee

## Overview

**Creative North Star: "The Roaster's Notebook"**

A working engineer's log kept on warm paper in espresso ink. The site reads like a notebook that lives next to the grinder: precise, honest, a little stained by use. Personality comes from the coffee identity and nowhere else: the palette is roasted, the dark mode is a darker roast rather than a black screen, and the brand mark is set in the same serif as the writing. Everything else is in service of reading long technical prose.

The home page is the notebook's cover: the brand and three photo doors set straight on the paper in the top-left corner, like a menu chalked onto the first page, outside the masthead-and-footer shell used by every other page.

Density is calm. Index pages are flat, year-grouped lists in the spirit of antfu.me; reading pages narrow to a 68-character measure. Hover is a quiet colour shift rather than a lift or a colour explosion. Controls stay quiet in muted sans until touched.

This system rejects three looks outright: generic SaaS (gradient heroes, glass cards, purple/blue tech palettes), the hacker terminal (neon on black, scanlines, matrix green), and the platform-default Medium/Substack template that carries no identity of its own.

**Key Characteristics:**
- Warm paper neutrals and espresso ink; a single cinnamon-roast accent.
- Dark mode is roasted brown, never neutral grey or black.
- Serif for what you read, sans for what you operate.
- Flat paper surfaces with no cards; one deep shadow reserved for the search dialog.
- Gently rounded corners, 44px touch targets, visible ember focus rings.
- Self-hosted fonts, no external runtime assets.
- One signature surface, the home menu, carries the site's framing; it sits on the page itself, never in a panel.

## Colors

A single-origin palette: oat and filter-paper neutrals, espresso and cold-brew inks, and one roasted accent that changes from cinnamon to caramel crema in the dark.

All colors are stored as space-separated RGB channels in CSS custom properties (`--bg: 250 246 238`) and applied as `rgb(var(--token))` or `rgb(var(--token) / alpha)`, so any token can take an alpha without a second definition. Component and frontmatter references use light-theme names; each has a dark counterpart listed below.

### Primary
- **Cinnamon Roast** (`--accent`, light): list markers, blockquote rules, link underlines at 55% alpha, card hover borders at 55%, selection at 28%. The accent appears as a line or a mark, rarely as a fill.
- **Cinnamon Roast Deep** (`--accent-strong`, light): link and brand text, active nav, table header text, tag emphasis. The accent color that carries text.
- **Caramel Crema / Caramel Crema Light** (`--accent` / `--accent-strong`, dark): the same two roles in dark themes.
- **Ember Focus / Caramel Focus** (`--focus`): the 3px focus outline, offset 3px. Never reused for decoration.

### Neutral
- **Oat Paper / Roasted Ground** (`--bg`): page background; also the inset field behind the search input.
- **Filter Paper / Dark Roast** (`--surface`): cards, tables, blockquotes, definition lists.
- **Fresh Filter / Medium Roast** (`--surface-raised`): the search dialog and skip link: the topmost sheet.
- **Espresso Ink / Cream** (`--text`): body text and headings.
- **Cold Brew / Milky Tan** (`--muted`): nav links, dates, summaries, hints, footer.
- **Kraft Line / Burnt Kraft** (`--border`): all hairlines, card borders, dividers, resting underlines.
- **Latte Foam / Dark Latte Foam** (`--chip-bg`) with **Mocha Label / Dark Mocha Label** (`--chip-text`): table headers, spoilers, inline math.
- **Parchment Code / Grounds** (`--code-bg`) with **Code Ink / Dark Code Ink** (`--code-text`): inline code and code blocks.
- **Roast Shadow** (`--shadow-color`): warm brown shadow base in light themes; pure black in dark themes.

### Menu Board (home only)
The home menu has no color family of its own: the title uses Cinnamon Roast Deep (Caramel Crema Light in the dark), door names use Espresso Ink (Cream), leaders and routes use Cold Brew (Milky Tan), and focus uses the site focus color. It draws no hairlines. Scene plates fall back to Latte Foam while a photo loads.

A functional sub-palette (`--syntax-*` for Chroma highlighting, `--task-*` for Neorg TODO/DOING/DONE/CANCELLED states) lives in `internal/web/styles/input.css`. It is confined to code blocks and task lists and never appears in site chrome.

### Named Rules
**The Unboxed Menu Rule.** The home menu is never put in a panel, card, or block, and draws no rules or hairlines. The title and doors sit directly on the page, anchored top-left with leaders running to the right margin, and the scene plates are square, borderless, and shadowless, in the spirit of borretti.me.

**The One Roast Rule.** Cinnamon Roast (Caramel Crema in the dark) is the only accent hue in the interface. A second chrome accent dilutes the identity; syntax and task colors stay inside content.

**The Roasted, Not Black Rule.** Dark surfaces are browns with warm undertones (`27 20 16` and up). Neutral greys and pure black are prohibited on surfaces.

## Typography

**Editorial Font:** Source Serif 4 (with Georgia, serif)
**UI Font:** Source Sans 3 (with ui-sans-serif, system-ui, sans-serif)
**Code Font:** system monospace stack (ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas)

**Character:** A book serif for reading paired with a humanist sans for labels: the notebook's handwriting and its printed index tabs. Both are self-hosted as WOFF2 at 400 and 600 only, and the regular cuts are preloaded.

### Hierarchy
- **Brand** (600, 1.05rem): the "danielscoffee" masthead mark, in the serif and the deep accent.
- **Headline** (600, 2.25rem, 1.2, -0.02em): page and article titles.
- **Title** (600, 1.5rem, 1.3, -0.02em): section headings on index pages.
- **Board Title / Board Row** (600, clamp(2rem → 2.5rem) at 1.1 for the title; clamp(1.5rem → 2rem) at 1.15 for door names and the intro line): the home board only, in the serif.
- **Card Title** (600, clamp(1.3rem, 3vw, 1.55rem), 1.25, -0.015em): post and project card titles, in the serif.
- **Body Prose** (400, clamp(1.03rem → 1.12rem), 1.82): article bodies, in the serif, at a 68ch measure. In-article h2/h3/h4 scale between 1.15rem and 1.9rem.
- **Body** (400, 1rem, 1.6): UI copy, summaries, page subtitles.
- **Label** (500 requested, 0.82rem): nav links, toggles, social links. Only 400 and 600 are shipped, so 500 resolves to the 400 cut.
- **Mono** (0.875rem, 1.5): code blocks; inline code at 0.92em.

**Observed drift:** headline and title currently inherit the UI sans because `.page-title`, `.post-title`, and `.section-title` set no family, although the 2026-07-12 refresh intended serif headings. Resolve deliberately; do not copy the drift into new surfaces without deciding.

### Named Rules
**The Ink and Label Rule.** The serif carries what you read (prose, card titles, the brand); the sans carries what you operate (navigation, metadata, controls, dates).

**The Measure Rule.** Running prose never exceeds 68ch.

## Layout

Two shells: a **wide shell** (`min(100% - 2rem, 72rem)`) for the masthead, footer, and index grids, and a **reading measure** (68ch) for intros and prose. Index lists (blog, projects, tag views) are flat and year-grouped: a large outlined year, then one row per entry, 1rem apart, with 2.5rem between years. Entries start at about 0.75 of the year's size, so its outline only grazes the first row. Project devlogs are a flat list without year marks: each row shows the full date ("Oct 10, 2026") and reading time.

The home page has its own shell: no masthead, no footer. The menu spans the full page width inside `clamp(1.5rem, 5vw, 3rem)` of page padding (48px on desktop), anchored top-left, so the dotted leaders run to the right margin and a growing plate pushes only the rows below it. At ≤480px the dotted leader and route label hide.

Vertical rhythm: main content pads 2.5rem top and 3rem bottom; sections are separated by 2.5rem; cards pad with `clamp(1.25rem, 3vw, 1.75rem)`.

Responsive behavior: inner pages use one masthead row at every width: the wordmark and trail at left; theme and search at right. Rather than wrap, the trail shortens only as far as it must to share the wordmark's row: full labels (`projects / blog engine`), then short (`proj / blog`), then tiny (`p / b`); screen readers always hear the full labels. At ≤560px the `Ctrl+K` hint is hidden and the shell gutter tightens to 0.625rem per side; at ≤400px the theme button shows its mode alone ("Light", aria-label still "Theme: Light") and the slashes tighten; at ≤340px the wordmark steps down to 1.15rem. Without JavaScript, narrow screens start on short labels and the actions wrap under the row if needed. At 480px the search dialog goes nearly full-screen.

## Elevation & Depth

Flat. Surfaces are distinguished by tone (Oat Paper → Filter Paper → Fresh Filter) and Kraft Line hairlines, never by shadow. The one shadow in the system belongs to the search dialog, warm and tinted with Roast Shadow, never grey.

### Shadow Vocabulary
- **Home menu** (`box-shadow: none`): no panel and no shadow; it is part of the page (The Unboxed Menu Rule).
- **Dialog** (`box-shadow: 0 1.5rem 4rem rgb(var(--shadow-color) / 0.28)`): the search dialog only, over a 58% Roast Shadow backdrop.

### Named Rules
**The Flat Page Rule.** Nothing on the page casts a shadow or lifts on hover: lists, menus, and headers sit on the paper. Only the search dialog floats above it.

## Shapes

Gently curved, never pill-shaped, never sharp. Controls, chips, inputs, and inline code use the small radius (0.375rem); cards and the dialog use the medium radius (0.75rem). Code blocks and definition panels use 0.5rem, content images 0.65rem, inline code and spoilers 0.25rem. Borders are always 1px Kraft Line; the thick line is the 0.25rem Cinnamon Roast rule on blockquotes. The home menu's scene plates are the deliberate exception to the radius scale: square, borderless, and flat (The Unboxed Menu Rule).

### Named Rules
**The Gentle Corner Rule.** Two radii do the work: 6px for things you press, 12px for things you read. The home menu's scene plates are the one square exception.

## Components

Warm and tactile: controls feel like paper and ink you can touch, with restraint.

### Buttons
- **Shape:** borderless ghost buttons for masthead actions (theme toggle, search trigger); outlined small-radius square (44×44px) for the dialog close.
- **Ghost:** Cold Brew label text, transparent, 0.35rem × 0.55rem padding, 44px min height.
- **Hover / Focus:** text shifts to Cinnamon Roast Deep over 180ms; focus adds the 3px Ember Focus ring.
- **Search trigger:** carries a `Ctrl+K` key hint in a hairline-bordered small-radius box at 0.65rem.

### Index List
Flat and year-grouped, after antfu.me/posts. No cards, borders, summaries, or tag chips.
- **Year mark:** the year as an `<h2>`, Source Serif 4 600 at clamp(4.5rem, 12vw, 8rem), transparent fill with a 1.5px Cold Brew outline at 60% alpha, positioned behind its group.
- **Entry:** the title link in the editorial serif (400, clamp(1.1rem, 2.4vw, 1.25rem)) in Espresso Ink, turning Cinnamon Roast Deep on hover and focus; then the meta in muted sans at 0.875rem: "Oct 10 · 4 min" for posts (reading time at 200 words a minute from the source), the date alone for projects. Rows wrap the meta under the title on narrow screens.
- **Page intro:** none visible. Blog and Projects keep a screen-reader `<h1>`; tag views show "Tagged with <tag>".

### Article Meta
Under an article title: the date, then each tag as a plain muted link separated by "·" (posts link to `/tag/<tag>`, projects to `/projects?tag=<tag>`). Hover and focus turn a tag Cinnamon Roast Deep with an underline. This is where tags live now that index lists carry none.

### Inputs / Fields
- **Style:** 48px tall, Oat Paper field set into the Fresh Filter dialog, 1px Kraft Line, small radius, 1rem text.
- **Focus:** border turns Cinnamon Roast with a 2px Ember Focus outline, offset 2px.
- **States:** search results show loading (accent text), empty (dashed Kraft Line box, centered muted text), and error (accent border) messages.

### Navigation
- **Masthead (inner pages):** the `danielscoffee` wordmark linking home, a smaller breadcrumb trail beside it, then theme and search at the far right, at every width, inside the wide shell. The home board is the site's navigation, so the header carries no sibling section links. Transparent and borderless; the footer carries the only top border and the GitHub, LinkedIn, and RSS links. Pages without a section (404, 500) show the wordmark alone. The home page has neither masthead nor footer.
- **Section index pages:** Blog, Projects, and About show no visible page title, because the section label already names them. Each keeps a screen-reader-only `<h1>`; Blog and Projects open with their one-line subtitle, About with its first paragraph. Posts, projects, devlogs, and tag pages keep their visible titles.
- **Section label:** the section name in lowercase, in the Brand style (Source Serif 4 600, 1.05rem, Cinnamon Roast Deep), 44px tall, linking to the section index; hover and focus turn it Espresso Ink.
- **Breadcrumb:** the wordmark (brand serif 600, `clamp(1.25rem, 1.1rem + 0.5vw, 1.45rem)`, Espresso Ink, turning Cinnamon Roast Deep on hover) is the home crumb; the trail after it is a `nav` labelled "Breadcrumb" with an ordered list, lowercase serif at 0.95rem, baseline-aligned with the wordmark. It is how readers go back: every crumb except the current page links up the site, never back through history. The last crumb, the page one step up from here, is Cinnamon Roast Deep at 600 (Espresso Ink on hover); earlier crumbs are Cold Brew, turning Cinnamon Roast Deep on hover and focus (180ms). On a section's own index the section crumb is the current page (`aria-current="page"`), not a link. Posts and tag pages read `danielscoffee / blog`, projects `danielscoffee / projects`, devlogs add their project. Each crumb carries full, short, and tiny labels; `crumbs-fit.js` picks the longest that fits (see Responsive behavior). Separators are muted slashes hidden from screen readers; every link keeps a 44px target. Articles carry no back link above the body.
- **Theme button:** a two-state Light/Dark switch labelled "Theme: Light" or "Theme: Dark". With no saved choice the site follows the OS (`data-theme="system"`) and the label shows the theme it resolves to; a click saves the opposite theme. The switch crossfades the whole page over 700ms on `cubic-bezier(0.45, 0, 0.55, 1)` (a typed same-document view transition, `theme`), so screen brightness ramps instead of flipping; browsers without typed view transitions ease colours over the same 700ms instead. The fade moves nothing, so it stays under reduced motion.
- **Section transition:** pages opt in to cross-document view transitions. Each door name on `/` shares a `view-transition-name` (`section-blog`, `section-projects`, `section-about`) with the section crumb on that section's pages, so "Blog" glides into the header (450ms, `cubic-bezier(0.2, 0, 0, 1)`) while the page crossfades (220ms), and back again on the home crumb. Door names with no partner fade at the page's 220ms instead of lingering. Reduced motion turns view transitions off; browsers without support simply navigate.
- **Links:** Cold Brew 0.82rem labels with a 44px hit area; hover and focus turn them Cinnamon Roast Deep. No underline, no pill highlight.
- **Footer:** brand name left, RSS/GitHub/LinkedIn right, muted 0.85rem, 44px targets.

### Menu Board (signature)
The site's doors, set straight on the page from the top-left corner across the full width (no panel, no hairlines; see The Unboxed Menu Rule). Anatomy: the serif brand title in the accent, then a list of rows separated only by space. Each row is a 44px-minimum link: scene plate, serif door name, a dotted Cold Brew leader running to the right page margin, and the route (`/blog`) in muted sans at the far right. Hover and focus darken the route and leader to the text color; focus draws the site's 3px focus outline offset 4px.

### Scene Plate (signature)
A square, borderless window onto a 400×400 photograph of Daniel's coffee rituals (writing at the table, pour-over at home, the café takeaway). Photos are CC0 from StockSnap, self-hosted as JPEGs under 60KB, credited in `internal/web/assets/img/scenes/CREDITS.md`, with the source embedded in each file. The image always renders at 192px, centered and cropped by the frame, so the frame opens like a shutter and never rescales the photo.
- **Sizes:** 96px at rest (shows the central 100-300 of the photo), 128px static on touch (67-333), 192px open.
- **Interaction:** on hover-capable devices the frame grows 96px → 192px over 0.3s ease-out on hover and focus-visible. Growth is held while the intro plays (`data-intro="play"`). Touch devices get static 128px frames with no transition.
- **Use:** plates frame navigation surfaces such as the home board. They are not content images and do not appear in article bodies or cards.

### Favicon
`internal/web/assets/favicon.svg`: an Oat Paper cup with Caramel Crema steam on an Espresso Ink tile (7px radius on a 32px grid). Strokes stay at 2px or more on that grid so the cup still reads at 16px. It is linked from the shared head as an SVG icon.

### Code Block
A Parchment Code frame (1px Kraft Line, 0.5rem radius) holding a caption bar over the highlighted code.
- **Bar:** the language mark (16px, monochrome `currentColor`, from Simple Icons, CC0) and the language name in the Label style, in Cold Brew. Unknown languages show a drawn `</>` mark and the name as written; no language shows "Code".
- **Copy button:** ghost style, 44px tall, a 1.5-stroke copy icon plus "Copy". It ships `hidden` and `code-copy.js` reveals it, so readers without JavaScript never see a dead control. A click copies the code, swaps to a check and "Copied" in Cinnamon Roast Deep for 1.6s, and announces "Code copied" through a polite live region. If clipboard access fails, it selects the code instead.
- **Rule:** logos stay monochrome. Brand colours would add hues to site chrome (The One Roast Rule).

### Search Dialog (signature)
A command palette on the topmost Fresh Filter sheet: serif heading, labelled 48px input, filter hint (`blog`, `projects`), and a scrolling result list of small-radius Filter Paper rows with an uppercase type label, title, and summary. It opens with Ctrl/⌘+K, closes on Esc, and returns focus to the trigger.

## Do's and Don'ts

### Do:
- **Do** apply color as `rgb(var(--token))` or `rgb(var(--token) / alpha)`; never hard-code a channel value in a component.
- **Do** define every new token for light, dark, and system-dark at once.
- **Do** keep every interactive target at 44px minimum and keep the 3px Ember Focus outline visible.
- **Do** keep prose at the 68ch measure and in Source Serif 4.
- **Do** use 180ms ease for color, border, background, and lift transitions, and let `prefers-reduced-motion` collapse them.
- **Do** keep the theme switch as a 700ms whole-page crossfade on `cubic-bezier(0.45, 0, 0.55, 1)`; it protects eyes from a full-screen brightness flip, and it is the second named exception to 180ms.
- **Do** treat the scene-plate frame-grow (0.3s ease-out) as the first named exception to 180ms; it is a size change that needs the longer travel. Any other motion longer than 180ms needs the same explicit justification.
- **Do** keep the board intro single-shot: the intro line holds for about 4.6s and fades out (350ms), then the three rows settle upward over 400ms with an 80ms stagger, finishing near 5.2s; it plays once per tab session, any key, click, or focus skips it, and `prefers-reduced-motion` skips it entirely.
- **Do** crop scene photos for the frame: the subject sits inside the central 100-300 window of the 400×400 crop so the 96px rest frame shows it whole, and the outer ring adds the context revealed on grow. Use only CC0 or equivalent photos, self-hosted, under 60KB, with the source embedded and listed in CREDITS.md.
- **Do** keep the inner-page header to the breadcrumb, theme, and search; section navigation lives on the home board.
- **Do** edit `internal/web/styles/input.css` and `.templ` sources; `output.css` is generated.
- **Do** keep coffee motifs in the site's framing (home, masthead, footer, empty and error states), where they say who wrote this, and keep them out of article bodies.

### Don't:
- **Don't** use generic SaaS styling: gradient heroes, glassmorphism cards, purple/blue tech palettes, centered marketing blocks.
- **Don't** use a hacker-terminal look: neon on black, monospace everything, scanlines, matrix green.
- **Don't** settle for a platform-default Medium/Substack reading template with no identity of its own.
- **Don't** introduce a second accent hue in site chrome (The One Roast Rule).
- **Don't** use neutral grey or pure black surfaces in dark mode (The Roasted, Not Black Rule).
- **Don't** load fonts, images, or scripts from external origins at runtime, or add inline scripts; the CSP forbids it.
- **Don't** use scene plates or the menu board outside framing surfaces such as the home page.
- **Don't** add shadows, cards, or hover lifts to page content; only the search dialog floats (The Flat Page Rule).
