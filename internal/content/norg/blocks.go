package norg

import (
	"fmt"
	"strconv"
	"strings"
)

func collectBlockLines(lines []string, start int, end string) ([]string, int, bool) {
	block := make([]string, 0)
	for i := start; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		if strings.TrimSpace(line) == end {
			return block, i + 1, true
		}
		block = append(block, line)
	}
	return block, len(lines), false
}

func parseNorgBlocks(lines []string) ([]norgNode, error) {
	nodes := make([]norgNode, 0)
	pendingAttrs := map[string]string{}
	for i := 0; i < len(lines); {
		line := strings.TrimRight(lines[i], "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			i++
			continue
		}

		if key, value, ok := parseCarryoverMeta(trimmed); ok {
			pendingAttrs[key] = value
			i++
			continue
		}

		attrs := takeAttrs(pendingAttrs)

		if lang, ok := parseAtCodeStart(trimmed); ok {
			codeLines, next, closed := collectBlockLines(lines, i+1, "@end")
			if !closed {
				return nil, fmt.Errorf("unclosed @code block")
			}
			i = next
			nodes = append(nodes, norgNode{kind: norgCode, lang: lang, code: strings.Join(codeLines, "\n"), attrs: attrs})
			continue
		}

		if isFenceStart(trimmed) {
			lang := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			codeLines, next, closed := collectBlockLines(lines, i+1, "```")
			if !closed {
				return nil, fmt.Errorf("unclosed code fence")
			}
			i = next
			nodes = append(nodes, norgNode{kind: norgCode, lang: lang, code: strings.Join(codeLines, "\n"), attrs: attrs})
			continue
		}

		if isTableStart(trimmed) {
			tableLines, next, closed := collectBlockLines(lines, i+1, "@end")
			if !closed {
				return nil, fmt.Errorf("unclosed @table block")
			}
			i = next
			tableHTML, err := parseMarkdownWrapperTable(tableLines)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, norgNode{kind: norgTable, tableHTML: tableHTML, attrs: attrs})
			continue
		}

		if text, ok := parseQuoteLine(trimmed); ok {
			quoteLines := []string{text}
			i++
			for i < len(lines) {
				next := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
				parsed, qok := parseQuoteLine(next)
				if !qok {
					break
				}
				quoteLines = append(quoteLines, parsed)
				i++
			}
			nodes = append(nodes, norgNode{kind: norgQuote, text: strings.Join(quoteLines, " "), attrs: attrs})
			continue
		}

		if item, ok := parseDefinitionLine(trimmed); ok {
			defs := []norgDefinitionItem{item}
			i++
			for i < len(lines) {
				next := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
				parsed, dok := parseDefinitionLine(next)
				if !dok {
					break
				}
				defs = append(defs, parsed)
				i++
			}
			nodes = append(nodes, norgNode{kind: norgDefinitionList, definitionItems: defs, attrs: attrs})
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
			nodes = append(nodes, norgNode{kind: norgTaskList, taskItems: items, attrs: attrs})
			continue
		}

		if src, ok := parseDotImageLine(trimmed); ok {
			alt := imageAltFromSrc(src)
			nodes = append(nodes, norgNode{kind: norgParagraph, text: fmt.Sprintf("![%s](%s)", alt, src), attrs: attrs})
			i++
			continue
		}

		if level, text, ok := parseHeadingLine(trimmed); ok {
			nodes = append(nodes, norgNode{kind: norgHeading, level: level, text: text, attrs: attrs})
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
			nodes = append(nodes, norgNode{kind: norgUL, items: items, attrs: attrs})
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
			nodes = append(nodes, norgNode{kind: norgOL, items: items, attrs: attrs})
			continue
		}

		paragraphLines := []string{trimmed}
		i++
		for i < len(lines) {
			next := strings.TrimSpace(strings.TrimRight(lines[i], "\r"))
			if _, ok := parseAtCodeStart(next); next == "" || ok || isTableStart(next) || isFenceStart(next) || isTaskLine(next) || isDotImageLine(next) || isHeadingLine(next) || isUnorderedLine(next) || isOrderedLine(next) {
				break
			}
			if _, qok := parseQuoteLine(next); qok {
				break
			}
			if _, dok := parseDefinitionLine(next); dok {
				break
			}
			if _, _, mok := parseCarryoverMeta(next); mok {
				break
			}
			paragraphLines = append(paragraphLines, next)
			i++
		}
		nodes = append(nodes, norgNode{kind: norgParagraph, text: strings.Join(paragraphLines, " "), attrs: attrs})
	}
	return nodes, nil
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

func isTableStart(line string) bool {
	return strings.TrimSpace(line) == "@table"
}

func parseAtCodeStart(line string) (string, bool) {
	if !strings.HasPrefix(line, "@code") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "@code")), true
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
