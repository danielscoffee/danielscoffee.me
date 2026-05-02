(() => {
  const key = "theme-preference";
  const modes = ["system", "light", "dark"];
  const labels = {
    system: "System",
    light: "Light",
    dark: "Dark",
  };

  const root = document.documentElement;
  const button = document.getElementById("theme-toggle");
  if (!button) return;

  const normalize = (value) => (modes.includes(value) ? value : "system");

  const readPreference = () => {
    try {
      return normalize(localStorage.getItem(key));
    } catch {
      return normalize(root.getAttribute("data-theme"));
    }
  };

  const writePreference = (value) => {
    try {
      localStorage.setItem(key, value);
    } catch {
      // Ignore storage failures (private mode, etc.)
    }
  };

  const apply = (value) => {
    const mode = normalize(value);
    root.setAttribute("data-theme", mode);
    const label = `Theme: ${labels[mode]}`;
    button.textContent = label;
    button.setAttribute("aria-label", label);
    button.dataset.themeMode = mode;
  };

  const cycle = (value) => {
    const idx = modes.indexOf(normalize(value));
    return modes[(idx + 1) % modes.length];
  };

  apply(readPreference());

  button.addEventListener("click", () => {
    const next = cycle(readPreference());
    writePreference(next);
    apply(next);
  });

  const media = window.matchMedia("(prefers-color-scheme: dark)");
  const handleSchemeChange = () => {
    if (readPreference() === "system") apply("system");
  };

  if (typeof media.addEventListener === "function") {
    media.addEventListener("change", handleSchemeChange);
  } else if (typeof media.addListener === "function") {
    media.addListener(handleSchemeChange);
  }
})();
