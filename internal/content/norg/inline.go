package norg

import (
	"fmt"
	"html"
	"maps"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	norgImagePattern = regexp.MustCompile(`!\{([^}]+)\}\[([^\]]*)\]|\{([^}]+)\}\[([^\]]*)\]\(image\)`)
	inlinePattern    = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)|\[([^\]]+)\]\(([^)]+)\)`)
	carryoverPattern = regexp.MustCompile(`^#([a-zA-Z0-9_-]+)\s+(.+)$`)
)

func renderInline(text string) (string, error) {
	text = normalizeNorgImageInline(text)
	idxs := inlinePattern.FindAllStringSubmatchIndex(text, -1)
	if len(idxs) == 0 {
		return renderTextWithModifiers(text), nil
	}

	var b strings.Builder
	cursor := 0
	for _, idx := range idxs {
		start := idx[0]
		end := idx[1]
		b.WriteString(renderTextWithModifiers(text[cursor:start]))

		if idx[2] != -1 && idx[4] != -1 {
			alt := html.EscapeString(text[idx[2]:idx[3]])
			rawSrc := text[idx[4]:idx[5]]
			src, err := validateCDNImageURL(rawSrc)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, `<img class="post-image" src="%s" alt="%s" loading="lazy" decoding="async"/>`, html.EscapeString(src), alt)
		} else {
			label := renderTextWithModifiers(text[idx[6]:idx[7]])
			href, err := validateLinkURL(text[idx[8]:idx[9]])
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, `<a href="%s">%s</a>`, html.EscapeString(href), label)
		}
		cursor = end
	}
	b.WriteString(renderTextWithModifiers(text[cursor:]))
	return b.String(), nil
}

func renderTextWithModifiers(raw string) string {
	var b strings.Builder
	for i := 0; i < len(raw); {
		delim := raw[i]
		if !isModifierDelimiter(delim) || isEscaped(raw, i) {
			if delim == '\\' && i+1 < len(raw) && isModifierDelimiter(raw[i+1]) {
				b.WriteString(html.EscapeString(string(raw[i+1])))
				i += 2
				continue
			}
			r, size := utf8.DecodeRuneInString(raw[i:])
			b.WriteString(html.EscapeString(string(r)))
			i += size
			continue
		}

		j := findClosingModifier(raw, i+1, delim)
		if j == -1 {
			b.WriteString(html.EscapeString(string(delim)))
			i++
			continue
		}

		inner := raw[i+1 : j]
		switch delim {
		case '*':
			fmt.Fprintf(&b, "<strong>%s</strong>", renderTextWithModifiers(inner))
		case '/':
			fmt.Fprintf(&b, "<em>%s</em>", renderTextWithModifiers(inner))
		case '_':
			fmt.Fprintf(&b, "<u>%s</u>", renderTextWithModifiers(inner))
		case '!':
			fmt.Fprintf(&b, "<span class=\"spoiler\">%s</span>", renderTextWithModifiers(inner))
		case '$':
			fmt.Fprintf(&b, "<span class=\"math-latex\">%s</span>", html.EscapeString(inner))
		}
		i = j + 1
	}
	return b.String()
}

func isModifierDelimiter(ch byte) bool {
	return strings.ContainsRune("*/_!$", rune(ch))
}

func isEscaped(s string, index int) bool {
	if index <= 0 {
		return false
	}
	count := 0
	for i := index - 1; i >= 0 && s[i] == '\\'; i-- {
		count++
	}
	return count%2 == 1
}

func findClosingModifier(s string, start int, delim byte) int {
	for i := start; i < len(s); i++ {
		if s[i] == '\n' {
			return -1
		}
		if s[i] == delim && !isEscaped(s, i) {
			if i == start {
				return -1
			}
			return i
		}
	}
	return -1
}

func parseCarryoverMeta(line string) (string, string, bool) {
	m := carryoverPattern.FindStringSubmatch(line)
	if len(m) != 3 {
		return "", "", false
	}
	return strings.ToLower(strings.TrimSpace(m[1])), strings.TrimSpace(m[2]), true
}

func takeAttrs(attrs map[string]string) map[string]string {
	if len(attrs) == 0 {
		return nil
	}
	out := maps.Clone(attrs)
	clear(attrs)
	return out
}

func renderAttrs(attrs map[string]string) string {
	var b strings.Builder
	for _, k := range slices.Sorted(maps.Keys(attrs)) {
		fmt.Fprintf(&b, " data-%s=\"%s\"", html.EscapeString(k), html.EscapeString(attrs[k]))
	}
	return b.String()
}

func parseQuoteLine(line string) (string, bool) {
	if !strings.HasPrefix(line, ">") {
		return "", false
	}
	text := strings.TrimSpace(line)
	for strings.HasPrefix(text, ">") {
		text = strings.TrimSpace(strings.TrimPrefix(text, ">"))
	}
	return text, true
}

func parseDefinitionLine(line string) (norgDefinitionItem, bool) {
	if !strings.HasPrefix(line, "$") {
		return norgDefinitionItem{}, false
	}
	raw := strings.TrimSpace(strings.TrimPrefix(line, "$"))
	if raw == "" {
		return norgDefinitionItem{}, false
	}
	if term, body, ok := strings.Cut(raw, ":"); ok {
		if strings.TrimSpace(body) != "" {
			return norgDefinitionItem{term: strings.TrimSpace(term), body: strings.TrimSpace(body)}, true
		}
	}
	return norgDefinitionItem{body: raw}, true
}

func validateLinkURL(raw string) (string, error) {
	for _, r := range raw {
		if r < ' ' || r == '\x7f' || r == '\\' {
			return "", fmt.Errorf("invalid link url %q", raw)
		}
	}

	href := strings.TrimSpace(raw)
	u, err := url.Parse(href)
	if err != nil || href == "" || (u.Host != "" && u.Scheme == "") {
		return "", fmt.Errorf("invalid link url %q", raw)
	}

	switch strings.ToLower(u.Scheme) {
	case "":
		if !strings.HasPrefix(href, "/") &&
			!strings.HasPrefix(href, "./") &&
			!strings.HasPrefix(href, "../") &&
			!strings.HasPrefix(href, "#") &&
			!strings.HasPrefix(href, "?") {
			return "", fmt.Errorf("invalid link url %q", raw)
		}
	case "http", "https":
		if u.Host == "" {
			return "", fmt.Errorf("invalid link url %q", raw)
		}
	case "mailto":
		if u.Opaque == "" {
			return "", fmt.Errorf("invalid link url %q", raw)
		}
	default:
		return "", fmt.Errorf("invalid link url %q", raw)
	}

	return u.String(), nil
}

func validateCDNImageURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("invalid image url %q", raw)
	}
	if u.Scheme != "https" || u.Host == "" {
		return "", fmt.Errorf("image url must be https CDN url: %q", raw)
	}

	host := strings.ToLower(u.Hostname())
	if !strings.Contains(host, "cdn") {
		return "", fmt.Errorf("image host must be CDN: %q", raw)
	}

	cleanPath := path.Clean("/" + strings.TrimPrefix(u.Path, "/"))
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("image path traversal blocked: %q", raw)
	}

	u.Path = cleanPath
	u.RawPath = ""
	return u.String(), nil
}

func parseDotImageLine(line string) (string, bool) {
	if !strings.HasPrefix(line, ".image ") {
		return "", false
	}
	src := strings.TrimSpace(strings.TrimPrefix(line, ".image"))
	if src == "" {
		return "", false
	}
	return src, true
}

func isDotImageLine(line string) bool {
	_, ok := parseDotImageLine(strings.TrimSpace(line))
	return ok
}

func imageAltFromSrc(src string) string {
	clean := strings.TrimSpace(src)
	if clean == "" {
		return "image"
	}
	parts := strings.Split(clean, "/")
	name := parts[len(parts)-1]
	if dot := strings.LastIndex(name, "."); dot > 0 {
		name = name[:dot]
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "image"
	}
	return name
}

func normalizeNorgImageInline(text string) string {
	idxs := norgImagePattern.FindAllStringSubmatchIndex(text, -1)
	if len(idxs) == 0 {
		return text
	}

	var b strings.Builder
	cursor := 0
	for _, idx := range idxs {
		start := idx[0]
		end := idx[1]
		b.WriteString(text[cursor:start])

		src := ""
		alt := ""
		if idx[2] != -1 {
			src = strings.TrimSpace(text[idx[2]:idx[3]])
			alt = strings.TrimSpace(text[idx[4]:idx[5]])
		} else {
			src = strings.TrimSpace(text[idx[6]:idx[7]])
			alt = strings.TrimSpace(text[idx[8]:idx[9]])
		}
		if alt == "" {
			alt = imageAltFromSrc(src)
		}
		fmt.Fprintf(&b, "![%s](%s)", alt, src)
		cursor = end
	}
	b.WriteString(text[cursor:])
	return b.String()
}
