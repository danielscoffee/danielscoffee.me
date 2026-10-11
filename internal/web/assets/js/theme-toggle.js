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
    button.textContent = label;
    button.setAttribute("aria-label", label);
    button.dataset.themeMode = mode;
  };

  render();

  button.addEventListener("click", () => {
    preference = effective() === "dark" ? "light" : "dark";
    writePreference(preference);
    render();
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
