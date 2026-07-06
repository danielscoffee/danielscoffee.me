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

	const resultLink = (href) => {
		const link = document.createElement("a");
		link.setAttribute(
			"href",
			typeof href === "string" && href.startsWith("/") ? href : "#",
		);
		return link;
	};

	const messageItem = (message) => {
		const li = document.createElement("li");
		li.className = "search-result";

		const link = resultLink("#");
		const title = document.createElement("div");
		title.className = "search-result-title";
		title.textContent = message;

		link.appendChild(title);
		li.appendChild(link);
		return li;
	};

	const render = (items) => {
		results.replaceChildren();
		if (!items.length) {
			results.appendChild(messageItem("No results"));
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

	let timer = null;
	const search = () => {
		const q = input.value.trim();
		if (!q) {
			results.replaceChildren();
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
				results.replaceChildren(messageItem("Search unavailable"));
			}
		}, 120);
	};

	trigger.addEventListener("click", open);
	input.addEventListener("input", search);

	modal.addEventListener("click", (event) => {
		const rect = modal.getBoundingClientRect();
		const inside =
			rect.top <= event.clientY &&
			event.clientY <= rect.top + rect.height &&
			rect.left <= event.clientX &&
			event.clientX <= rect.left + rect.width;
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
