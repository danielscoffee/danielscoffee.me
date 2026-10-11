// Copy buttons on code blocks. The parser ships them hidden; this reveals
// them, copies the block's text, and falls back to selecting it.
(() => {
  const buttons = document.querySelectorAll(".code-copy");
  if (!buttons.length) return;

  const status = document.createElement("span");
  status.className = "sr-only";
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  document.body.append(status);

  const selectText = (node) => {
    const range = document.createRange();
    range.selectNodeContents(node);
    const selection = window.getSelection();
    selection.removeAllRanges();
    selection.addRange(range);
  };

  for (const button of buttons) {
    const pre = button.closest(".code-block")?.querySelector("pre");
    const label = button.querySelector(".code-copy-text");
    if (!pre || !label) continue;

    let timer = 0;
    const show = (state, text, announcement) => {
      button.dataset.state = state;
      label.textContent = text;
      status.textContent = announcement;
      clearTimeout(timer);
      timer = setTimeout(() => {
        delete button.dataset.state;
        label.textContent = "Copy";
      }, 1600);
    };

    button.hidden = false;
    button.addEventListener("click", async () => {
      const text = pre.textContent.replace(/\n$/, "");
      try {
        await navigator.clipboard.writeText(text);
        show("copied", "Copied", "Code copied");
      } catch {
        selectText(pre);
        show("selected", "Selected", "Code selected, press Ctrl+C or Command+C to copy");
      }
    });
  }
})();
