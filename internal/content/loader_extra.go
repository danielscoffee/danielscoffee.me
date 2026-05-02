package content

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
)

func LoadProjects(dir string) ([]Project, error) {
	entries, err := loadEntries(dir)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0, len(entries))
	for _, entry := range entries {
		if entry.meta.Draft {
			continue
		}
		projects = append(projects, Project{
			Title:    entry.meta.Title,
			Slug:     entry.meta.Slug,
			Date:     entry.meta.Date,
			Summary:  entry.meta.Summary,
			Tags:     entry.meta.Tags,
			Draft:    entry.meta.Draft,
			BodyMD:   entry.body,
			BodyHTML: template.HTML(entry.htmlBody),
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Date > projects[j].Date
	})

	return projects, nil
}

func LoadPage(path string) (Page, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Page{}, fmt.Errorf("read %s: %w", path, err)
	}

	meta, body, preRenderedHTML, err := splitFrontMatter(string(raw), filepath.Ext(path))
	if err != nil {
		return Page{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if meta.Draft {
		return Page{}, fmt.Errorf("page %s is marked draft", path)
	}

	htmlBody := preRenderedHTML
	if htmlBody == "" {
		htmlBody, err = renderMarkdown(body)
		if err != nil {
			return Page{}, fmt.Errorf("render markdown %s: %w", path, err)
		}
	}

	return Page{
		Title:    meta.Title,
		Slug:     meta.Slug,
		Date:     meta.Date,
		Summary:  meta.Summary,
		BodyMD:   body,
		BodyHTML: template.HTML(htmlBody),
	}, nil
}

type contentEntry struct {
	meta     frontMatter
	body     string
	htmlBody string
}

func loadEntries(dir string) ([]contentEntry, error) {
	mdFiles, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}
	norgFiles, err := filepath.Glob(filepath.Join(dir, "*.norg"))
	if err != nil {
		return nil, err
	}
	files := append(mdFiles, norgFiles...)

	entries := make([]contentEntry, 0, len(files))
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}

		meta, body, preRenderedHTML, err := splitFrontMatter(string(raw), filepath.Ext(file))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}

		htmlBody := preRenderedHTML
		if htmlBody == "" {
			htmlBody, err = renderMarkdown(body)
			if err != nil {
				return nil, fmt.Errorf("render markdown %s: %w", file, err)
			}
		}

		entries = append(entries, contentEntry{meta: meta, body: body, htmlBody: htmlBody})
	}
	return entries, nil
}
