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
  board-espresso: "rgb(54 35 25)"
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
  board-roast: "rgb(78 54 39)"
typography:
  brand:
    fontFamily: "Source Serif 4, Georgia, serif"
    fontSize: "1.05rem"
    fontWeight: 600
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
  card-padding: "clamp(1.25rem, 3vw, 1.75rem)"
  grid-gap: "1.25rem"
  section-gap: "2.5rem"
  touch-target: "44px"
  content-reading: "68ch"
  content-wide: "72rem"
  board-width: "40rem"
  plate-rest: "96px"
  plate-touch: "128px"
  plate-open: "192px"
components:
  brand:
    textColor: "{colors.cinnamon-roast-deep}"
    typography: "{typography.brand}"
    height: "44px"
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
  card-editorial:
    backgroundColor: "{colors.filter-paper}"
    rounded: "{rounded.md}"
    padding: "{spacing.card-padding}"
  card-title-link:
    textColor: "{colors.cinnamon-roast-deep}"
    typography: "{typography.card-title}"
  card-title-link-hover:
    textColor: "{colors.espresso-ink}"
  tag-chip:
    backgroundColor: "{colors.latte-foam}"
    textColor: "{colors.mocha-label}"
    rounded: "{rounded.sm}"
    padding: "0.25rem 0.5rem"
  tag-chip-hover:
    textColor: "{colors.cinnamon-roast-deep}"
  input-search:
    backgroundColor: "{colors.oat-paper}"
    textColor: "{colors.espresso-ink}"
    rounded: "{rounded.sm}"
    height: "48px"
    padding: "0.7rem 0.85rem"
  menu-board:
    backgroundColor: "{colors.board-espresso}"
    textColor: "{colors.cream}"
    padding: "clamp(1.25rem, 4vw, 2.5rem)"
    width: "40rem"
  menu-board-dark:
    backgroundColor: "{colors.board-roast}"
    textColor: "{colors.cream}"
  board-row:
    textColor: "{colors.cream}"
    typography: "{typography.board-row}"
    height: "44px"
  scene-plate:
    backgroundColor: "{colors.kraft-frame}"
    size: "96px"
  site-menu-toggle:
    textColor: "{colors.cold-brew}"
    typography: "{typography.label}"
    rounded: "{rounded.sm}"
    height: "44px"
  site-menu-toggle-hover:
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

The home page is the notebook's cover: a single espresso menu board with three framed doors, set outside the masthead-and-footer shell used by every other page.

Density is calm. Index pages are a wide shell of paper cards; reading pages narrow to a 68-character measure. Components are warm and tactile: cards are sheets of filter paper with a faint warm shadow, chips are dabs of latte foam, and hover is a gentle physical lift rather than a color explosion. Controls stay quiet in muted sans until touched.

This system rejects three looks outright: generic SaaS (gradient heroes, glass cards, purple/blue tech palettes), the hacker terminal (neon on black, scanlines, matrix green), and the platform-default Medium/Substack template that carries no identity of its own.

**Key Characteristics:**
- Warm paper neutrals and espresso ink; a single cinnamon-roast accent.
- Dark mode is roasted brown, never neutral grey or black.
- Serif for what you read, sans for what you operate.
- Flat paper surfaces, softly lifted; one deep shadow reserved for the search dialog.
- Gently rounded corners, 44px touch targets, visible ember focus rings.
- Self-hosted fonts, no external runtime assets.
- One signature surface, the menu board, carries the site's framing; inner pages stay paper.

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
- **Latte Foam / Dark Latte Foam** (`--chip-bg`) with **Mocha Label / Dark Mocha Label** (`--chip-text`): tag chips, table headers, spoilers, inline math.
- **Parchment Code / Grounds** (`--code-bg`) with **Code Ink / Dark Code Ink** (`--code-text`): inline code and code blocks.
- **Roast Shadow** (`--shadow-color`): warm brown shadow base in light themes; pure black in dark themes.

### Menu Board (home only)
The board is a dark surface on both themes; its tokens (`--board*`) are defined in light, dark, and system-dark, and only `.menu-board` and its children consume them.
- **Board Espresso** (`--board`, light): the board fill, identical in value to Espresso Ink.
- **Board Roast** (`--board`, dark): a lifted medium-roast brown, lighter than the Dark Roast page surfaces so the board reads as an object on the page.
- **Cream / Milky Tan** (`--board-text` / `--board-muted`): board text and its hairlines, leaders, and route labels (at 40-60% alpha); same values in both themes.
- **Kraft Frame** (`--board-frame`): the scene-plate backing fill, visible only while a photo loads; same value in both themes.
- **Caramel Focus** (`--board-focus`): the board's focus outline and selection tint; same value in both themes.

A functional sub-palette (`--syntax-*` for Chroma highlighting, `--task-*` for Neorg TODO/DOING/DONE/CANCELLED states) lives in `internal/web/styles/input.css`. It is confined to code blocks and task lists and never appears in site chrome.

### Named Rules
**The Board Stays Dark Rule.** The menu board is espresso or board-roast in every theme and never takes the paper surface colors. It is the only large dark surface in light mode, which is why it reads as a board and not a card.

**The Flat Board Rule.** The menu board and its scene plates are square-cornered, borderless, and shadowless, in the spirit of borretti.me. Tone alone separates board from page and plate from board.

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

Two shells: a **wide shell** (`min(100% - 2rem, 72rem)`) for the masthead, footer, and index grids, and a **reading measure** (68ch) for intros and prose. Index lists are an auto-fit grid of cards, `minmax(min(100%, 25rem), 1fr)` with a 1.25rem gap. Project devlogs use a fixed one/two/three-column grid at 640px and 900px, where the first devlog spans two columns and two rows.

The home page has its own shell: no masthead, no footer. The menu board is a 40rem-max centered column, top-anchored near the resting board's vertical center so a growing plate pushes only the rows below it. The resting half-height is tracked per mode in `--board-rest-half` (16rem default; 14.75rem at ≤480px; 19rem on touch; 17.5rem on touch at ≤480px). At ≤480px the dotted leader and route label hide and board padding tightens to 1.25rem × 1rem.

Vertical rhythm: main content pads 2.5rem top and 3rem bottom; sections are separated by 2.5rem; cards pad with `clamp(1.25rem, 3vw, 1.75rem)`.

Responsive behavior: inner pages use one masthead at every width: brand mark left, a `<details>` "Menu" disclosure at right. Opening it reveals one panel with nav, social links, theme, and search, laid out as a right-aligned wrapping row under a Kraft Line top border; at ≤560px the panel becomes a single column of full-width rows and the shell gutter tightens to 0.625rem per side. On touch devices the `Ctrl+K` hint is hidden. At 480px the search dialog goes nearly full-screen.

## Elevation & Depth

Paper, softly lifted. Surfaces are flat sheets distinguished first by tone (Oat Paper → Filter Paper → Fresh Filter) and a Kraft Line hairline; shadows are faint, warm, and tinted with Roast Shadow, never grey.

### Shadow Vocabulary
- **Resting sheet** (`box-shadow: 0 0.25rem 0.9rem rgb(var(--shadow-color) / 0.06)`): cards and devlog tiles at rest.
- **Lifted sheet** (`box-shadow: 0 0.4rem 1.2rem rgb(var(--shadow-color) / 0.09)`): card hover and focus-within, paired with `translateY(-1px)` on editorial cards.
- **Board** (`box-shadow: none`): the menu board is flat; it is made heavy by tone alone (The Flat Board Rule).
- **Dialog** (`box-shadow: 0 1.5rem 4rem rgb(var(--shadow-color) / 0.28)`): the search dialog only, over a 58% Roast Shadow backdrop.

### Named Rules
**The Soft Lift Rule.** At rest, a sheet barely casts a shadow. Interaction lifts it by one pixel and a few percent of shadow. Only the search dialog floats.

## Shapes

Gently curved, never pill-shaped, never sharp. Controls, chips, inputs, and inline code use the small radius (0.375rem); cards and the dialog use the medium radius (0.75rem). Code blocks and definition panels use 0.5rem, content images 0.65rem, inline code and spoilers 0.25rem. Borders are always 1px Kraft Line; the thick line is the 0.25rem Cinnamon Roast rule on blockquotes. The menu board and scene plates are the deliberate exception to the radius scale: square, borderless, and flat (The Flat Board Rule).

### Named Rules
**The Gentle Corner Rule.** Two radii do the work: 6px for things you press, 12px for things you read. The home board and its plates are the one square exception.

## Components

Warm and tactile: controls feel like paper and ink you can touch, with restraint.

### Buttons
- **Shape:** borderless ghost buttons for masthead actions (theme toggle, search trigger); outlined small-radius square (44×44px) for the dialog close.
- **Ghost:** Cold Brew label text, transparent, 0.35rem × 0.55rem padding, 44px min height.
- **Hover / Focus:** text shifts to Cinnamon Roast Deep over 180ms; focus adds the 3px Ember Focus ring.
- **Search trigger:** carries a `Ctrl+K` key hint in a hairline-bordered small-radius box at 0.65rem.

### Chips
- **Style:** Latte Foam fill, Mocha Label text at 0.75rem, small radius, 0.25rem × 0.5rem padding, no border.
- **State:** hover/focus turns the text Cinnamon Roast Deep. Chips are links to tag pages.

### Cards / Containers
- **Corner Style:** medium radius (0.75rem).
- **Background:** Filter Paper.
- **Shadow Strategy:** Resting sheet; Lifted sheet plus 1px rise on hover (see Elevation).
- **Border:** 1px Kraft Line, shifting to Cinnamon Roast at 55% on hover.
- **Internal Padding:** clamp(1.25rem, 3vw, 1.75rem).
- **Anatomy:** serif title link (Cinnamon Roast Deep with a Kraft Line underline offset 0.2em, turning Espresso Ink with a currentColor underline on hover), muted summary at 60ch, then a metadata row pinned to the bottom with the date and tag chips.

### Inputs / Fields
- **Style:** 48px tall, Oat Paper field set into the Fresh Filter dialog, 1px Kraft Line, small radius, 1rem text.
- **Focus:** border turns Cinnamon Roast with a 2px Ember Focus outline, offset 2px.
- **States:** search results show loading (accent text), empty (dashed Kraft Line box, centered muted text), and error (accent border) messages.

### Navigation
- **Masthead (inner pages):** serif brand mark at left and one "Menu" disclosure at right, at every width, inside the wide shell. Transparent and borderless; the footer carries the only top border. The home page has neither masthead nor footer.
- **Menu toggle:** 44px summary in the Label style, small radius, with a drawn plus that collapses to a minus when open (two 1.5px bars, 180ms). Cold Brew at rest; Cinnamon Roast Deep on hover, focus, and open.
- **Menu panel:** nav links (Blog, Projects, About), then GitHub and LinkedIn, then the theme and search buttons, in one wrapping right-aligned row (single column ≤560px). The theme button's visible text is "Theme: <mode>".
- **Links:** Cold Brew 0.82rem labels with a 44px hit area; hover and focus turn them Cinnamon Roast Deep. No underline, no pill highlight.
- **Footer:** brand name left, RSS/GitHub/LinkedIn right, muted 0.85rem, 44px targets.

### Menu Board (signature)
A flat dark espresso panel (40rem max, square corners, no border, no shadow) that holds the site's doors. Anatomy: a serif board title over a muted hairline, then a list of rows separated by 40%-alpha Milky Tan hairlines. Each row is a 44px-minimum link: scene plate, serif door name, a dotted leader, and the route (`/blog`) in muted sans. Hover and focus brighten the route and leader to Cream; focus draws a 3px Caramel Focus outline offset 4px. The board never uses paper surface tokens (see The Board Stays Dark Rule).

### Scene Plate (signature)
A square, borderless window onto a 400×400 photograph of Daniel's coffee rituals (writing at the table, pour-over at home, the café takeaway). Photos are CC0 from StockSnap, self-hosted as JPEGs under 60KB, credited in `internal/web/assets/img/scenes/CREDITS.md`, with the source embedded in each file. The image always renders at 192px, centered and cropped by the frame, so the frame opens like a shutter and never rescales the photo.
- **Sizes:** 96px at rest (shows the central 100-300 of the photo), 128px static on touch (67-333), 192px open.
- **Interaction:** on hover-capable devices the frame grows 96px → 192px over 0.3s ease-out on hover and focus-visible. Growth is held while the intro plays (`data-intro="play"`). Touch devices get static 128px frames with no transition.
- **Use:** plates frame navigation surfaces such as the home board. They are not content images and do not appear in article bodies or cards.

### Search Dialog (signature)
A command palette on the topmost Fresh Filter sheet: serif heading, labelled 48px input, filter hint (`blog`, `projects`), and a scrolling result list of small-radius Filter Paper rows with an uppercase type label, title, and summary. It opens with Ctrl/⌘+K, closes on Esc, and returns focus to the trigger.

## Do's and Don'ts

### Do:
- **Do** apply color as `rgb(var(--token))` or `rgb(var(--token) / alpha)`; never hard-code a channel value in a component.
- **Do** define every new token for light, dark, and system-dark at once.
- **Do** keep every interactive target at 44px minimum and keep the 3px Ember Focus outline visible.
- **Do** keep prose at the 68ch measure and in Source Serif 4.
- **Do** use 180ms ease for color, border, background, and lift transitions, and let `prefers-reduced-motion` collapse them.
- **Do** treat the scene-plate frame-grow (0.3s ease-out) as the one named exception to 180ms; it is a size change that needs the longer travel. Any other motion longer than 180ms needs the same explicit justification.
- **Do** keep the board intro single-shot: the intro line holds for about 4.6s and fades out (350ms), then the three rows settle upward over 400ms with an 80ms stagger, finishing near 5.2s; it plays once per tab session, any key, click, or focus skips it, and `prefers-reduced-motion` skips it entirely.
- **Do** crop scene photos for the frame: the subject sits inside the central 100-300 window of the 400×400 crop so the 96px rest frame shows it whole, and the outer ring adds the context revealed on grow. Use only CC0 or equivalent photos, self-hosted, under 60KB, with the source embedded and listed in CREDITS.md.
- **Do** give the inner-page header one "Menu" disclosure at every width rather than a separate desktop nav.
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
- **Don't** cast heavy shadows on anything but the search dialog (The Soft Lift Rule).
