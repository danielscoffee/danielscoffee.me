(() => {
  const key = "theme-preference";
  const validModes = new Set(["system", "light", "dark"]);

  let mode = "system";
  try {
    const stored = localStorage.getItem(key);
    if (stored && validModes.has(stored)) mode = stored;
  } catch {
    mode = "system";
  }

  document.documentElement.setAttribute("data-theme", mode);
})();
