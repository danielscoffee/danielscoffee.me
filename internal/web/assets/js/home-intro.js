// Plays the home intro once per tab session; any key, click, or focus skips it.
(() => {
  const root = document.documentElement;
  const key = "home-intro-seen";
  const skip = () => {
    root.dataset.intro = "skip";
  };

  let seen = false;
  try {
    seen = Boolean(sessionStorage.getItem(key));
    if (!seen) sessionStorage.setItem(key, "1");
  } catch {
    // Storage unavailable: the CSS intro still plays and stays skippable.
  }

  if (seen || window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    skip();
    return;
  }

  // Frames stay shut until the rows settle (the CSS intro ends near 5.2s).
  root.dataset.intro = "play";
  setTimeout(() => {
    if (root.dataset.intro === "play") root.dataset.intro = "done";
  }, 5250);

  for (const type of ["keydown", "pointerdown", "click", "focusin"]) {
    document.addEventListener(type, skip, { once: true, capture: true });
  }
})();
