package norg

import "strings"

type norgNodeKind int

const (
	norgHeading norgNodeKind = iota
	norgParagraph
	norgUL
	norgOL
	norgTaskList
	norgCode
	norgTable
	norgQuote
	norgDefinitionList
)

type norgNode struct {
	kind            norgNodeKind
	level           int
	text            string
	items           []string
	taskItems       []norgTaskItem
	lang            string
	code            string
	tableHTML       string
	definitionItems []norgDefinitionItem
	attrs           map[string]string
}

type norgTaskItem struct {
	state string
	text  string
}

type norgDefinitionItem struct {
	term string
	body string
}

func Parse(raw string) (Meta, string, string, error) {
	meta, bodyLines, err := splitNorgFrontMatter(raw)
	if err != nil {
		return Meta{}, "", "", err
	}

	nodes, err := parseNorgBlocks(bodyLines)
	if err != nil {
		return Meta{}, "", "", err
	}

	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	htmlBody, err := renderNorgHTML(nodes)
	if err != nil {
		return Meta{}, "", "", err
	}
	return meta, body, htmlBody, nil
}
