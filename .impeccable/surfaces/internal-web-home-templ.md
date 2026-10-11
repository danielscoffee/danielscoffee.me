---
version: 1
slug: "internal-web-home-templ"
primary_target: "internal/web/home.templ"
related_targets: ["internal/web/base.templ","internal/web/styles/input.css"]
---

Scope: `/` home page (new `internal/web/home.templ`), plus the slimmed inner-page masthead in `base.templ`. Visitor mode: Experience. Full brief: `docs/plans/2026-10-10-home-coffee-nav-design.md`.

Audience: cold arrivals from all five PRODUCT.md audiences. Job: know whose place this is and pick a door (Blog, Projects, About) within seconds. Proof: scenes drawn from Daniel's real rituals (pour-over at home, café regular).

Memorable moment: the intro line on the board giving way to three doors whose frames grow to reveal the full coffee scene.

Intro timing: the line holds about 4.6s and the doors land near 5s (user request).

## Direction contract

THESIS: The front door is the café's menu board: three framed scenes on one dark-roast board, not a hero, a card grid, or a nav bar. It refuses the centred headline over a row of equal cards.

OWN-WORLD: Oat-paper page; one flat, square, borderless dark-roast board (lifted medium roast in dark theme) with cream Source Serif 4 lettering; square, borderless photo plates of Daniel's real coffee rituals (CC0, self-hosted); dotted leaders ending in small UI-type routes.

STORY: The visitor reads "Backend developer. Systems enjoyer. Production bug survivor." on the board, sees it give way to Blog, Projects, and About, hovers a frame to see the pour-over station or the café takeaway open up, and walks through a door.

FIRST VIEWPORT: Board centred on both axes at about 640px wide; header `danielscoffee` at about 2.5rem serif; three rows below, each a 96px frame, a name at about 2rem serif, a leader, and a route. No masthead or footer. The primary action is the three rows themselves.

FORM: Café menu board, position 5 of 7 on the ordered list, seed key 5a9d0deb. Signature interaction: borretti frame-grow 96px to 192px over 0.3s ease-out on hover and focus-visible, revealing a fixed 192px scene. Motion grammar: intro line fades out, rows settle with a short stagger, then only hover motion.

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
