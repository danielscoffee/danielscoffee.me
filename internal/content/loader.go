package content

import (
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
	"time"
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
	return loadPublished(dir, "post", func(entry contentEntry) Post {
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

func loadPublished[T publishedContent](dir, kind string, convert func(contentEntry) T) ([]T, error) {
	entries, err := loadEntries(dir)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(entries))
	seenSlugs := make(map[string]struct{})
	for _, entry := range entries {
		if entry.meta.Draft {
			continue
		}
		slug := normalizeKey(entry.meta.Slug)
		if _, dup := seenSlugs[slug]; dup {
			return nil, fmt.Errorf("duplicate %s slug %q", kind, entry.meta.Slug)
		}
		seenSlugs[slug] = struct{}{}
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
	if !slugPattern.MatchString(meta.Slug) {
		return fmt.Errorf("invalid slug %q: use lowercase letters, numbers, and single hyphens", meta.Slug)
	}
	if _, err := time.Parse("2006-01-02", meta.Date); err != nil {
		return fmt.Errorf("invalid date %q: use YYYY-MM-DD", meta.Date)
	}
	for _, tag := range meta.Tags {
		if !tagPattern.MatchString(tag) {
			return fmt.Errorf("invalid tag %q: use lowercase letters, numbers, and single hyphens", tag)
		}
	}
	return nil
}

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	tagPattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)
