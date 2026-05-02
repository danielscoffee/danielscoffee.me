// Ctrl+K search modal
(() => {
  const modal = document.getElementById("search-modal");
  const input = document.getElementById("search-input");
  const results = document.getElementById("search-results");
  const trigger = document.getElementById("search-trigger");
  if (!modal || !input || !results || !trigger) return;

  const open = () => {
    if (!modal.open) modal.showModal();
    input.focus();
    input.select();
  };

  const close = () => {
    if (modal.open) modal.close();
  };

  const render = (items) => {
    results.innerHTML = "";
    if (!items.length) {
      results.innerHTML = `<li class="search-result"><a href="#"><div class="search-result-title">No results</div></a></li>`;
      return;
    }

    for (const item of items) {
      const li = document.createElement("li");
      li.className = "search-result";
      li.innerHTML = `
        <a href="${item.url}">
          <div class="search-result-type">${item.type}</div>
          <div class="search-result-title">${item.title}</div>
          <div class="search-result-summary">${item.summary || ""}</div>
        </a>
      `;
      results.appendChild(li);
    }
  };

  let timer = null;
  const search = () => {
    const q = input.value.trim();
    if (!q) {
      results.innerHTML = "";
      return;
    }

    clearTimeout(timer);
    timer = setTimeout(async () => {
      try {
        const res = await fetch(`/search?q=${encodeURIComponent(q)}`);
        if (!res.ok) throw new Error("search failed");
        const data = await res.json();
        render(data.results || []);
      } catch {
        results.innerHTML = `<li class="search-result"><a href="#"><div class="search-result-title">Search unavailable</div></a></li>`;
      }
    }, 120);
  };

  trigger.addEventListener("click", open);
  input.addEventListener("input", search);

  modal.addEventListener("click", (event) => {
    const rect = modal.getBoundingClientRect();
    const inside = rect.top <= event.clientY && event.clientY <= rect.top + rect.height && rect.left <= event.clientX && event.clientX <= rect.left + rect.width;
    if (!inside) close();
  });

  document.addEventListener("keydown", (event) => {
    const key = event.key.toLowerCase();
    if ((event.ctrlKey || event.metaKey) && key === "k") {
      event.preventDefault();
      open();
      return;
    }

    if (key === "escape" && modal.open) {
      close();
    }
  });
})();
