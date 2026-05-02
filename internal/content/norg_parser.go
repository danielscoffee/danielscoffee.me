package content

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type norgNodeKind int

const (
	norgHeading norgNodeKind = iota
	norgParagraph
	norgUL
	norgOL
	norgTaskList
	norgCode
)

type norgNode struct {
	kind      norgNodeKind
	level     int
	text      string
	items     []string
	taskItems []norgTaskItem
	lang      string
	code      string
}

type norgTaskItem struct {
	state string
	text  string
}

func parseNorg(raw string) (frontMatter, string, string, error) {
	meta, bodyLines, err := splitNorgFrontMatter(raw)
	if err != nil {
		return frontMatter{}, "", "", err
	}

	nodes, err := parseNorgBlocks(bodyLines)
	if err != nil {
		return frontMatter{}, "", "", err
	}

	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	htmlBody := renderNorgHTML(nodes)
	return meta, body, htmlBody, nil
}

func splitNorgFrontMatter(raw string) (frontMatter, []string, error) {
	var meta frontMatter

	trimmed := strings.TrimSpace(raw)
	lines := strings.Split(trimmed, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "@document.meta" {
		return meta, nil, fmt.Errorf("missing frontmatter delimiter")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "@end" {
			end = i
			break
		}
	}
	if end == -1 {
		return meta, nil, fmt.Errorf("invalid frontmatter structure")
	}

	parsed, err := parseNeorgMeta(lines[1:end])
	if err != nil {
		return meta, nil, err
	}
	if err := validateFrontMatter(parsed); err != nil {
		return meta, nil, err
	}

	return parsed, lines[end+1:], nil
}

func parseNeorgMeta(lines []string) (frontMatter, error) {
	meta := frontMatter{}
	openListKey := ""

	for _, rawLine := range lines {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}

		if openListKey != "" {
			if strings.HasPrefix(trimmed, "-") {
				item := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
				if openListKey == "tags" && item != "" {
					meta.Tags = append(meta.Tags, stripQuotes(item))
					continue
				}
			}
			openListKey = ""
		}

		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			return meta, fmt.Errorf("invalid neorg metadata line %q", trimmed)
		}

		k := strings.ToLower(strings.TrimSpace(key))
		v := strings.TrimSpace(value)
		if v == "" {
			openListKey = k
			continue
		}

		switch k {
		case "title":
			meta.Title = stripQuotes(v)
		case "slug":
			meta.Slug = stripQuotes(v)
		case "date":
			meta.Date = stripQuotes(v)
		case "summary":
			meta.Summary = stripQuotes(v)
		case "draft":
			draft, err := strconv.ParseBool(stripQuotes(v))
			if err != nil {
				return meta, fmt.Errorf("invalid draft value %q", v)
			}
			meta.Draft = draft
		case "tags":
			meta.Tags = append(meta.Tags, parseTagValues(v)...)
		}
	}

	return meta, nil
}

func parseTagValues(raw string) []string {
	v := strings.TrimSpace(raw)
	if v == "" {
		return nil
	}

	if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
		v = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(v, "["), "]"))
	}

	parts := strings.Split(v, ",")
	if len(parts) == 1 {
		item := stripQuotes(strings.TrimSpace(parts[0]))
		if item == "" {
			return nil
		}
		return []string{item}
	}

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		item := stripQuotes(strings.TrimSpace(p))
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func stripQuotes(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 2 {
		if (strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"")) ||
			(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func parseNorgBlocks(lines []string) ([]norgNode, error) {
	nodes := make([]norgNode, 0)
	for i := 0; i < len(lines); {
		line := strings.TrimRight(lines[i], "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			i++
			continue
		}

		if isFenceStart(trimmed) {
			lang := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			i++
			codeLines := make([]string, 0)
			closed := false
			for i < len(lines) {
				current := strings.TrimRight(lines[i], "\r")
				if strings.TrimSpace(current) == "```" {
					closed = true
					i++
					break
				}
				codeLines = append(codeLines, current)
				i++
			}
			if !closed {
				return nil, fmt.Errorf("unclosed code fence")
			}
			nodes = append(nodes, norgNode{kind: norgCode, lang: lang, code: strings.Join(codeLines, "\n")})
			continue
		}

		if strings.HasPrefix(trimmed, "*** ") {
			state, text, ok := parseTaskLine(trimmed)
			if !ok {
				parts := strings.Fields(trimmed)
				if len(parts) >= 2 {
					return nil, fmt.Errorf("invalid task state %q", parts[1])
				}
				return nil, fmt.Errorf("invalid task state")
			}
			items := []norgTaskItem{{state: state, text: text}}
			i++
			for i < len(lines) {
				nextTrimmed := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
				if !strings.HasPrefix(nextTrimmed, "*** ") {
					break
				}
				nextState, nextText, nextOK := parseTaskLine(nextTrimmed)
				if !nextOK {
					nextParts := strings.Fields(nextTrimmed)
					if len(nextParts) >= 2 {
						return nil, fmt.Errorf("invalid task state %q", nextParts[1])
					}
					return nil, fmt.Errorf("invalid task state")
				}
				items = append(items, norgTaskItem{state: nextState, text: nextText})
				i++
			}
			nodes = append(nodes, norgNode{kind: norgTaskList, taskItems: items})
			continue
		}

		if level, text, ok := parseHeadingLine(trimmed); ok {
			nodes = append(nodes, norgNode{kind: norgHeading, level: level, text: text})
			i++
			continue
		}

		if item, ok := parseUnorderedLine(trimmed); ok {
			items := []string{item}
			i++
			for i < len(lines) {
				nextItem, nextOK := parseUnorderedLine(strings.TrimSpace(strings.TrimRight(lines[i], "\r")))
				if !nextOK {
					break
				}
				items = append(items, nextItem)
				i++
			}
			nodes = append(nodes, norgNode{kind: norgUL, items: items})
			continue
		}

		if item, ok := parseOrderedLine(trimmed); ok {
			items := []string{item}
			i++
			for i < len(lines) {
				nextItem, nextOK := parseOrderedLine(strings.TrimSpace(strings.TrimRight(lines[i], "\r")))
				if !nextOK {
					break
				}
				items = append(items, nextItem)
				i++
			}
			nodes = append(nodes, norgNode{kind: norgOL, items: items})
			continue
		}

		paragraphLines := []string{trimmed}
		i++
		for i < len(lines) {
			next := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
			if next == "" || isFenceStart(next) || isTaskLine(next) || isHeadingLine(next) || isUnorderedLine(next) || isOrderedLine(next) {
				break
			}
			paragraphLines = append(paragraphLines, next)
			i++
		}
		nodes = append(nodes, norgNode{kind: norgParagraph, text: strings.Join(paragraphLines, " ")})
	}
	return nodes, nil
}

func renderNorgHTML(nodes []norgNode) string {
	var b strings.Builder
	for _, n := range nodes {
		switch n.kind {
		case norgHeading:
			lvl := n.level
			if lvl < 1 {
				lvl = 1
			}
			if lvl > 6 {
				lvl = 6
			}
			fmt.Fprintf(&b, "<h%d>%s</h%d>\n", lvl, renderInline(n.text), lvl)
		case norgParagraph:
			fmt.Fprintf(&b, "<p>%s</p>\n", renderInline(n.text))
		case norgUL:
			b.WriteString("<ul>\n")
			for _, item := range n.items {
				fmt.Fprintf(&b, "<li>%s</li>\n", renderInline(item))
			}
			b.WriteString("</ul>\n")
		case norgOL:
			b.WriteString("<ol>\n")
			for _, item := range n.items {
				fmt.Fprintf(&b, "<li>%s</li>\n", renderInline(item))
			}
			b.WriteString("</ol>\n")
		case norgTaskList:
			b.WriteString("<ul class=\"task-list\">\n")
			for _, item := range n.taskItems {
				fmt.Fprintf(&b, "<li data-task-state=\"%s\">%s</li>\n", renderTaskState(item.state), renderInline(item.text))
			}
			b.WriteString("</ul>\n")
		case norgCode:
			b.WriteString(renderHighlightedCode(n.lang, n.code))
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func renderHighlightedCode(lang, code string) string {
	lexerName := strings.TrimSpace(lang)
	lexer := lexers.Get(lexerName)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return plainCodeHTML(lang, code)
	}

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var out bytes.Buffer
	if err := formatter.Format(&out, style, iterator); err != nil {
		return plainCodeHTML(lang, code)
	}

	rendered := strings.TrimSpace(out.String())
	rendered = chromaPreStylePattern.ReplaceAllString(rendered, `<pre class="chroma">`)
	return rendered
}

func plainCodeHTML(lang, code string) string {
	langClass := ""
	if lang != "" {
		langClass = fmt.Sprintf(" class=\"language-%s\"", html.EscapeString(lang))
	}
	return fmt.Sprintf("<pre><code%s>%s</code></pre>", langClass, html.EscapeString(code))
}

func renderTaskState(state string) string {
	switch state {
	case "TODO":
		return "todo"
	case "DOING":
		return "doing"
	case "DONE":
		return "done"
	case "CANCELLED":
		return "cancelled"
	default:
		return "todo"
	}
}

var (
	chromaPreStylePattern = regexp.MustCompile(`<pre style="[^"]*">`)
	linkPattern           = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
)

func renderInline(text string) string {
	idxs := linkPattern.FindAllStringSubmatchIndex(text, -1)
	if len(idxs) == 0 {
		return html.EscapeString(text)
	}

	var b strings.Builder
	cursor := 0
	for _, idx := range idxs {
		start := idx[0]
		end := idx[1]
		labelStart := idx[2]
		labelEnd := idx[3]
		hrefStart := idx[4]
		hrefEnd := idx[5]

		b.WriteString(html.EscapeString(text[cursor:start]))
		label := html.EscapeString(text[labelStart:labelEnd])
		href := html.EscapeString(text[hrefStart:hrefEnd])
		fmt.Fprintf(&b, `<a href="%s">%s</a>`, href, label)
		cursor = end
	}
	b.WriteString(html.EscapeString(text[cursor:]))
	return b.String()
}

func isFenceStart(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "```")
}

func parseTaskLine(line string) (string, string, bool) {
	if !strings.HasPrefix(line, "*** ") {
		return "", "", false
	}
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return "", "", false
	}
	state := parts[1]
	if _, ok := allowedTaskStates[state]; !ok {
		return "", "", false
	}
	return state, strings.Join(parts[2:], " "), true
}

func isTaskLine(line string) bool {
	_, _, ok := parseTaskLine(line)
	return ok
}

func parseHeadingLine(line string) (int, string, bool) {
	if !isHeadingLine(line) {
		return 0, "", false
	}
	level := 0
	for level < len(line) && line[level] == '*' {
		level++
	}
	return level, strings.TrimSpace(line[level:]), true
}

func isHeadingLine(line string) bool {
	if line == "" || line[0] != '*' {
		return false
	}
	i := 0
	for i < len(line) && line[i] == '*' {
		i++
	}
	return i < len(line) && line[i] == ' '
}

func parseUnorderedLine(line string) (string, bool) {
	if !strings.HasPrefix(line, "- ") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "- ")), true
}

func isUnorderedLine(line string) bool {
	_, ok := parseUnorderedLine(line)
	return ok
}

func parseOrderedLine(line string) (string, bool) {
	idx := strings.Index(line, ".")
	if idx <= 0 || idx+1 >= len(line) || line[idx+1] != ' ' {
		return "", false
	}
	if _, err := strconv.Atoi(line[:idx]); err != nil {
		return "", false
	}
	return strings.TrimSpace(line[idx+1:]), true
}

func isOrderedLine(line string) bool {
	_, ok := parseOrderedLine(line)
	return ok
}

var allowedTaskStates = map[string]struct{}{
	"TODO":      {},
	"DOING":     {},
	"DONE":      {},
	"CANCELLED": {},
}
