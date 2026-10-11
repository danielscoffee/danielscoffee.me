(() => {
  const key = "theme-preference";
  const labels = {
    light: "Light",
    dark: "Dark",
  };

  const root = document.documentElement;
  const button = document.getElementById("theme-toggle");
  if (!button) return;

  const media = window.matchMedia("(prefers-color-scheme: dark)");

  // A saved light/dark choice wins; with none, the page follows the OS
  // (data-theme="system") and the button shows what that resolves to.
  const readPreference = () => {
    try {
      const stored = localStorage.getItem(key);
      return stored === "light" || stored === "dark" ? stored : null;
    } catch {
      return null;
    }
  };

  const writePreference = (value) => {
    try {
      localStorage.setItem(key, value);
    } catch {
      // Ignore storage failures (private mode, etc.)
    }
  };

  let preference = readPreference();
  const effective = () => preference || (media.matches ? "dark" : "light");

  const render = () => {
    root.setAttribute("data-theme", preference || "system");
    const mode = effective();
    const label = `Theme: ${labels[mode]}`;
    // The prefix sits in its own span so narrow screens can show the mode alone.
    const prefix = document.createElement("span");
    prefix.className = "theme-prefix";
    prefix.textContent = "Theme:";
    button.replaceChildren(prefix, labels[mode]);
    button.setAttribute("aria-label", label);
    button.dataset.themeMode = mode;
  };

  render();

  const switchTheme = () => {
    preference = effective() === "dark" ? "light" : "dark";
    writePreference(preference);
    render();
  };

  // Switching crossfades the whole page over 700ms so screen brightness
  // ramps instead of flipping. Browsers with typed view transitions fade a
  // snapshot of the page; the rest ease their colors under .theme-fading.
  const fadeMs = 700;
  const typedTransitions =
    typeof document.startViewTransition === "function" &&
    typeof CSS !== "undefined" &&
    CSS.supports("selector(:active-view-transition-type(theme))");
  let fadeTimer = 0;

  button.addEventListener("click", () => {
    if (typedTransitions) {
      document.startViewTransition({ update: switchTheme, types: ["theme"] });
      return;
    }
    root.classList.add("theme-fading");
    switchTheme();
    clearTimeout(fadeTimer);
    fadeTimer = setTimeout(() => root.classList.remove("theme-fading"), fadeMs);
  });

  const handleSchemeChange = () => {
    if (!preference) render();
  };

  if (typeof media.addEventListener === "function") {
    media.addEventListener("change", handleSchemeChange);
  } else if (typeof media.addListener === "function") {
    media.addListener(handleSchemeChange);
  }
})();
