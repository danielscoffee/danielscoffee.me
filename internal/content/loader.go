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
			Published: publishedFromMeta(entry.meta),
			BodyMD:    entry.body,
			BodyHTML:  template.HTML(entry.htmlBody),
		}
	})
}

func publishedFromMeta(meta frontMatter) Published {
	return Published{
		Title:   meta.Title,
		Slug:    meta.Slug,
		Date:    meta.Date,
		Summary: meta.Summary,
		Tags:    meta.Tags,
		Draft:   meta.Draft,
	}
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
