package content

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func LoadProjects(dir string) ([]Project, error) {
	return loadPublished(dir, func(entry contentEntry) Project {
		return Project{
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

func LoadPage(path string) (Page, error) {
	entry, err := loadContentFile(path)
	if err != nil {
		return Page{}, err
	}

	if entry.meta.Draft {
		return Page{}, fmt.Errorf("page %s is marked draft", path)
	}

	return Page{
		Title:    entry.meta.Title,
		Slug:     entry.meta.Slug,
		Date:     entry.meta.Date,
		Summary:  entry.meta.Summary,
		BodyMD:   entry.body,
		BodyHTML: template.HTML(entry.htmlBody),
	}, nil
}

func loadContentFile(path string) (contentEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return contentEntry{}, fmt.Errorf("read %s: %w", path, err)
	}

	meta, body, preRenderedHTML, err := splitFrontMatter(string(raw), filepath.Ext(path))
	if err != nil {
		return contentEntry{}, fmt.Errorf("parse %s: %w", path, err)
	}

	if preRenderedHTML == "" {
		return contentEntry{}, fmt.Errorf("unsupported content format %q", filepath.Ext(path))
	}

	return contentEntry{meta: meta, body: body, htmlBody: preRenderedHTML}, nil
}

type contentEntry struct {
	meta     frontMatter
	body     string
	htmlBody string
}

func loadEntries(dir string) ([]contentEntry, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.norg"))
	if err != nil {
		return nil, err
	}

	entries := make([]contentEntry, 0, len(files))
	for _, file := range files {
		entry, err := loadContentFile(file)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
