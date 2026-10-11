// Shortens the masthead breadcrumb only as far as it must to stay on the
// wordmark's row: "projects / blog engine", then "proj / blog", then
// "p / b". Screen readers always hear the full labels.
(() => {
  const nav = document.querySelector(".site-crumbs");
  const row = nav?.closest(".masthead-inner");
  const brand = row?.querySelector(".site-brand");
  const actions = row?.querySelector(".header-actions");
  const crumbs = nav?.querySelectorAll(".crumb") ?? [];
  if (!row || !brand || !actions || !crumbs.length) return;

  const sameRow = (a, b) => {
    const ra = a.getBoundingClientRect();
    const rb = b.getBoundingClientRect();
    return ra.top < rb.bottom && rb.top < ra.bottom;
  };
  const fits = () =>
    sameRow(brand, nav) && sameRow(brand, actions) && sameRow(crumbs[0], crumbs[crumbs.length - 1]);

  const fit = () => {
    for (const size of ["full", "short", "tiny"]) {
      nav.dataset.fit = size;
      if (fits()) return;
    }
  };

  // Refit when the row's width changes, not when fitting changes its height.
  let width = -1;
  new ResizeObserver(([entry]) => {
    if (entry.contentRect.width === width) return;
    width = entry.contentRect.width;
    fit();
  }).observe(row);
  document.fonts?.ready.then(fit);
})();
