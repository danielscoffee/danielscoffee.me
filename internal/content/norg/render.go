package norg

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
	chromaPreStylePattern   = regexp.MustCompile(`<pre style="[^"]*">`)
	tableSeparatorCellRegex = regexp.MustCompile(`^:?-+:?$`)
)

func renderNorgHTML(nodes []norgNode) (string, error) {
	var b strings.Builder
	for _, n := range nodes {
		attrs := renderAttrs(n.attrs)
		switch n.kind {
		case norgHeading:
			lvl := min(max(n.level+1, 2), 6)
			text, err := renderInline(n.text)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "<h%d%s>%s</h%d>\n", lvl, attrs, text, lvl)
		case norgParagraph:
			text, err := renderInline(n.text)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "<p%s>%s</p>\n", attrs, text)
		case norgUL, norgOL:
			tag := "ul"
			if n.kind == norgOL {
				tag = "ol"
			}
			fmt.Fprintf(&b, "<%s%s>\n", tag, attrs)
			for _, item := range n.items {
				text, err := renderInline(item)
				if err != nil {
					return "", err
				}
				fmt.Fprintf(&b, "<li>%s</li>\n", text)
			}
			fmt.Fprintf(&b, "</%s>\n", tag)
		case norgTaskList:
			fmt.Fprintf(&b, "<ul class=\"task-list\"%s>\n", attrs)
			for _, item := range n.taskItems {
				text, err := renderInline(item.text)
				if err != nil {
					return "", err
				}
				fmt.Fprintf(&b, "<li data-task-state=\"%s\">%s</li>\n", strings.ToLower(item.state), text)
			}
			b.WriteString("</ul>\n")
		case norgCode, norgTable:
			body := n.tableHTML
			if n.kind == norgCode {
				body = renderHighlightedCode(n.lang, n.code)
			}
			if attrs != "" {
				fmt.Fprintf(&b, "<div%s>", attrs)
			}
			b.WriteString(body)
			if attrs != "" {
				b.WriteString("</div>")
			}
			b.WriteString("\n")
		case norgQuote:
			text, err := renderInline(n.text)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "<blockquote%s><p>%s</p></blockquote>\n", attrs, text)
		case norgDefinitionList:
			fmt.Fprintf(&b, "<dl class=\"norg-definitions\"%s>\n", attrs)
			for _, d := range n.definitionItems {
				if d.term != "" {
					term, err := renderInline(d.term)
					if err != nil {
						return "", err
					}
					fmt.Fprintf(&b, "<dt>%s</dt>\n", term)
				}
				body, err := renderInline(d.body)
				if err != nil {
					return "", err
				}
				fmt.Fprintf(&b, "<dd>%s</dd>\n", body)
			}
			b.WriteString("</dl>\n")
		}
	}
	return strings.TrimSpace(b.String()), nil
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
	rendered = strings.Replace(rendered, `<pre class="chroma">`, `<pre class="chroma" tabindex="0">`, 1)
	rendered = chromaPreStylePattern.ReplaceAllString(rendered, `<pre class="chroma" tabindex="0">`)
	return rendered
}

func plainCodeHTML(lang, code string) string {
	langClass := ""
	if lang != "" {
		langClass = fmt.Sprintf(" class=\"language-%s\"", html.EscapeString(lang))
	}
	return fmt.Sprintf("<pre tabindex=\"0\"><code%s>%s</code></pre>", langClass, html.EscapeString(code))
}

func parseMarkdownWrapperTable(lines []string) (string, error) {
	rows := make([][]string, 0, len(lines))
	for _, raw := range lines {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		cells, ok := parseTableRow(trimmed)
		if !ok {
			return "", fmt.Errorf("invalid @table row %q", raw)
		}
		rows = append(rows, cells)
	}
	if len(rows) < 2 {
		return "", fmt.Errorf("invalid @table structure")
	}

	headers := rows[0]
	sep := rows[1]
	if len(headers) == 0 || len(headers) != len(sep) {
		return "", fmt.Errorf("invalid @table structure")
	}
	for _, cell := range sep {
		if !tableSeparatorCellRegex.MatchString(strings.TrimSpace(cell)) {
			return "", fmt.Errorf("invalid @table separator")
		}
	}

	var b strings.Builder
	b.WriteString("<table tabindex=\"0\">\n<thead><tr>")
	for _, header := range headers {
		text, err := renderInline(header)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "<th>%s</th>", text)
	}
	b.WriteString("</tr></thead>\n<tbody>\n")
	for _, row := range rows[2:] {
		if len(row) != len(headers) {
			return "", fmt.Errorf("invalid @table row width")
		}
		b.WriteString("<tr>")
		for _, cell := range row {
			text, err := renderInline(cell)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "<td>%s</td>", text)
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("</tbody>\n</table>")
	return b.String(), nil
}

func parseTableRow(line string) ([]string, bool) {
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil, false
	}
	parts := strings.Split(line, "|")
	if len(parts) < 3 {
		return nil, false
	}
	cells := make([]string, 0, len(parts)-2)
	for _, part := range parts[1 : len(parts)-1] {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells, true
}
