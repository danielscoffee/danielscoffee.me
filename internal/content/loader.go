package content

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/danielscoffee/danielscoffee.me/internal/content/norg"
)

func LoadPosts(dir string) ([]Post, error) {
	return loadPublished(dir, "post", func(entry contentEntry) Post {
		return Post{
			Published: publishedFromMeta(entry.meta),
			BodyMD:    entry.body,
			BodyHTML:  template.HTML(entry.htmlBody),
		}
	})
}

func publishedFromMeta(meta norg.Meta) Published {
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

func splitFrontMatter(raw, ext string) (norg.Meta, string, string, error) {
	if strings.EqualFold(ext, ".norg") {
		return norg.Parse(raw)
	}
	return norg.Meta{}, "", "", fmt.Errorf("unsupported content format %q", ext)
}
