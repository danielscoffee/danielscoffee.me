// Ctrl+K / Cmd+K search modal
(() => {
	const modal = document.getElementById("search-modal");
	const input = document.getElementById("search-input");
	const results = document.getElementById("search-results");
	const trigger = document.getElementById("search-trigger");
	if (!modal || !input || !results || !trigger) return;

	let timer = null;
	let controller = null;

	const cancelSearch = () => {
		clearTimeout(timer);
		timer = null;
		if (controller) controller.abort();
		controller = null;
		results.removeAttribute("aria-busy");
	};

	const open = () => {
		if (!modal.open) modal.showModal();
		input.focus();
		input.select();
	};

	const close = () => {
		if (modal.open) modal.close();
	};

	const resultLink = (href) => {
		const link = document.createElement("a");
		link.setAttribute(
			"href",
			typeof href === "string" && /^\/(?!\/)/.test(href) && !href.includes("\\")
				? href
				: "#",
		);
		return link;
	};

	const messageClasses = {
		loading: "search-result-message-loading",
		empty: "search-result-message-empty",
		error: "search-result-message-error",
	};

	const messageItem = (message, state) => {
		const li = document.createElement("li");
		li.className = `search-result search-result-message ${messageClasses[state] || ""}`;
		li.textContent = message;
		return li;
	};

	const showMessage = (message, state) => {
		results.replaceChildren(messageItem(message, state));
	};

	const render = (items) => {
		results.replaceChildren();
		if (!items.length) {
			showMessage("No results", "empty");
			return;
		}

		for (const item of items) {
			const li = document.createElement("li");
			li.className = "search-result";
			const link = resultLink(item.url);
			const type = document.createElement("div");
			type.className = "search-result-type";
			type.textContent = item.type || "";
			const title = document.createElement("div");
			title.className = "search-result-title";
			title.textContent = item.title || "";
			const summary = document.createElement("div");
			summary.className = "search-result-summary";
			summary.textContent = item.summary || "";
			link.append(type, title, summary);
			li.appendChild(link);
			results.appendChild(li);
		}
	};

	const search = () => {
		cancelSearch();
		const q = input.value.trim();
		if (!q) {
			results.replaceChildren();
			return;
		}

		timer = setTimeout(async () => {
			controller = new AbortController();
			const request = controller;
			results.setAttribute("aria-busy", "true");
			showMessage("Loading…", "loading");
			try {
				const res = await fetch(`/search?q=${encodeURIComponent(q)}`, {
					signal: request.signal,
				});
				if (!res.ok) throw new Error("search failed");
				const data = await res.json();
				if (controller === request) render(data.results || []);
			} catch (error) {
				if (error.name !== "AbortError" && controller === request) {
					showMessage("Search unavailable", "error");
				}
			} finally {
				if (controller === request) {
					controller = null;
					results.removeAttribute("aria-busy");
				}
			}
		}, 120);
	};

	trigger.addEventListener("click", open);
	input.addEventListener("input", search);
	modal.addEventListener("close", () => {
		cancelSearch();
		input.value = "";
		results.replaceChildren();
		trigger.focus();
	});
	modal.addEventListener("click", (event) => {
		const rect = modal.getBoundingClientRect();
		const inside = rect.top <= event.clientY && event.clientY <= rect.bottom && rect.left <= event.clientX && event.clientX <= rect.right;
		if (!inside) close();
	});

	document.addEventListener("keydown", (event) => {
		const key = event.key.toLowerCase();
		if ((event.ctrlKey || event.metaKey) && key === "k") {
			event.preventDefault();
			open();
			return;
		}
		if (key === "escape" && modal.open) close();
	});
})();
