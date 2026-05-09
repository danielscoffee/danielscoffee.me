package content

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
)

type frontMatter struct {
	Title   string   `yaml:"title"`
	Slug    string   `yaml:"slug"`
	Date    string   `yaml:"date"`
	Summary string   `yaml:"summary"`
	Tags    []string `yaml:"tags"`
	Draft   bool     `yaml:"draft"`
}

func LoadPosts(dir string) ([]Post, error) {
	return loadPublished(dir, func(entry contentEntry) Post {
		return Post{
			Published: Published{
				Title:   entry.meta.Title,
				Slug:    entry.meta.Slug,
				Date:    entry.meta.Date,
				Summary: entry.meta.Summary,
				Tags:    entry.meta.Tags,
				Draft:   entry.meta.Draft,
			},
			BodyMD:   entry.body,
			BodyHTML: template.HTML(entry.htmlBody),
		}
	})
}

type publishedContent interface {
	publishedDate() string
}

func loadPublished[T publishedContent](dir string, convert func(contentEntry) T) ([]T, error) {
	entries, err := loadEntries(dir)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(entries))
	for _, entry := range entries {
		if entry.meta.Draft {
			continue
		}
		items = append(items, convert(entry))
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].publishedDate() > items[j].publishedDate()
	})

	return items, nil
}

func splitFrontMatter(raw, ext string) (frontMatter, string, string, error) {
	switch strings.ToLower(ext) {
	case ".norg":
		meta, body, html, err := parseNorg(raw)
		return meta, body, html, err
	default:
		return frontMatter{}, "", "", fmt.Errorf("unsupported content format %q", ext)
	}
}

func validateFrontMatter(meta frontMatter) error {
	if meta.Title == "" || meta.Slug == "" || meta.Date == "" {
		return fmt.Errorf("title, slug, and date are required")
	}
	return nil
}

