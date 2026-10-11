# Home Coffee Navigation Design

## Goal

Replace the `/` → `/blog` redirect with a centred navigation page in the spirit of borretti.me: three doors (Blog, Projects, About), each a coffee scene whose frame grows on hover. Slim the inner-page masthead to a brand mark plus a compact menu. Product context lives in `PRODUCT.md`; the visual system in `DESIGN.md`.

## Audience and outcome

Cold visitors landing on `/`: peers, recruiters, collaborators, clients, and Daniel himself. Visitor mode is Experience: the scenes lead and the interface recedes. Within seconds a visitor knows whose place this is and picks a door. The scenes are true to Daniel's rituals (pour-over at home, café regular), not generic café décor.

## Direction: The Menu Board

- A dark-roast café menu board centred on the oat-paper page. Board header: `danielscoffee` in Source Serif 4, cream.
- Three rows, ruled to one module: kraft-framed scene, large serif name, dotted leader, route in small UI type (`/blog`, `/projects`, `/about`).
- Signature interaction (borretti mechanic, exact): the frame grows 96px → 192px over 0.3s ease-out; the image inside is fixed at 192px and centred, with overflow hidden, so growth reveals more of the scene. Triggered by `:hover` and `:focus-visible`. 300ms ease-out is a deliberate exception to DESIGN.md's 180ms default.
- Flat, in the spirit of borretti.me: square corners, no border, no shadow on the board or its frames. Light theme: board in espresso brown; dark theme: board in a lifted medium roast. Lettering cream in both.

## Scenes

Three CC0 photographs from StockSnap, self-hosted at `internal/web/assets/img/scenes/` as 400×400 JPEG crops (under 60KB each) and displayed at 192px. Each subject sits inside the central 50% so it reads at the 96px crop; the outer ring adds context revealed on grow. Credits are in `CREDITS.md` beside the files, and each JPEG embeds its source. Decorative: `aria-hidden` on the frame, with the label text naming the link. (Revised 2026-10-10: the first build used hand-drawn SVGs; the user asked for real photos.)

- **Blog: café table.** A hand writing in an open notebook at a wooden table, with a cup and croissant revealed on grow.
- **Projects: pour-over station.** A dripper on a carafe on a scale, centred; the gooseneck kettle pouring and the window light revealed on grow.
- **About: the café takeaway.** A paper cup of latte on a café table, centred, with its lid set beside it.

## Intro sequence

- On first load the board shows only the intro line "Backend developer. Systems enjoyer. Production bug survivor." for about 4.6s, then the line fades out as the three rows settle in with a short stagger. Total about 5s (raised from about 2s at the user's request).
- The doors are in the DOM from first paint and clickable during the intro. Focus, any key, or a click skips it.
- CSS-driven, so it plays without JavaScript. A tiny external script (theme-init.js pattern, no inline script) marks the session so later visits to `/` in the same session show the doors immediately.
- `prefers-reduced-motion: reduce` skips the intro and shows the doors at once.

## Page structure

- **Home:** no masthead, no footer, only the board. `/` renders the home page with status 200; any other unmatched path still 404s.
- **Inner pages:** brand mark (links to `/`) plus a compact disclosure menu at every width containing Blog, Projects, About, GitHub, LinkedIn, the theme toggle, and search. Ctrl+K search keeps working. Footer unchanged.

## Responsive and accessibility

- Below about 480px the leader and route hide; each row is scene plus name.
- On `(hover: none)` devices, frames rest at about 128px and do not grow; tap navigates.
- 44px minimum targets, visible focus ring on the board (ember/caramel focus), 200% zoom reflows, light/dark/system themes.

## Constraints

- CSP: no inline scripts; no external fonts, images, or scripts.
- Edit `.templ` sources and `internal/web/styles/input.css`; `output.css` and `*_templ.go` are generated (`make generate`).
- No changes to other routes, content schema, search API, or article layout.
- Update the redirect test (`internal/http/blog_routes_test.go:98`), the masthead assertion (`:189`), and README routes.

## Anti-goals

Generic SaaS, hacker terminal, platform-template look, stock-café clichés.

## Validation

`make generate`, `go test ./...`, `go vet ./...`. Review `/` and an inner page at 375px, 768px, and 1440px in light and dark themes; check keyboard flow, focus visibility, the intro skip, reduced motion, and 200% zoom. Run the Impeccable detector on changed UI files.
