// The arrow above a post, project, or devlog returns to the previous page
// when the reader came from this site. Arriving from elsewhere, or opening
// it in a new tab, follows its href to the parent listing instead.
(() => {
  const link = document.querySelector(".page-back");
  if (!link || history.length < 2) return;

  let fromThisSite = false;
  try {
    fromThisSite = new URL(document.referrer).origin === location.origin;
  } catch {
    return;
  }
  if (!fromThisSite) return;

  link.addEventListener("click", (event) => {
    if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
    event.preventDefault();
    history.back();
  });
})();
