package norg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Meta struct {
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Draft   bool
}

func splitNorgFrontMatter(raw string) (Meta, []string, error) {
	var meta Meta

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
	if err := validateMeta(parsed); err != nil {
		return meta, nil, err
	}

	return parsed, lines[end+1:], nil
}

func parseNeorgMeta(lines []string) (Meta, error) {
	meta := Meta{}
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
	if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
		v = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(v, "["), "]"))
	}

	var out []string
	for _, p := range strings.Split(v, ",") {
		item := stripQuotes(strings.TrimSpace(p))
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func stripQuotes(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 2 && ((v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'')) {
		return v[1 : len(v)-1]
	}
	return v
}

func validateMeta(meta Meta) error {
	if meta.Title == "" || meta.Slug == "" || meta.Date == "" {
		return fmt.Errorf("title, slug, and date are required")
	}
	if !keyPattern.MatchString(meta.Slug) {
		return fmt.Errorf("invalid slug %q: use lowercase letters, numbers, and single hyphens", meta.Slug)
	}
	if _, err := time.Parse("2006-01-02", meta.Date); err != nil {
		return fmt.Errorf("invalid date %q: use YYYY-MM-DD", meta.Date)
	}
	for _, tag := range meta.Tags {
		if !keyPattern.MatchString(tag) {
			return fmt.Errorf("invalid tag %q: use lowercase letters, numbers, and single hyphens", tag)
		}
	}
	return nil
}

var keyPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
