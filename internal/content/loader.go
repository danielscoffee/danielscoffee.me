package content

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
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
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(files))
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}

		meta, md, err := splitFrontMatter(string(raw))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		if meta.Draft {
			continue
		}

		htmlBody, err := renderMarkdown(md)
		if err != nil {
			return nil, fmt.Errorf("render markdown %s: %w", file, err)
		}

		posts = append(posts, Post{
			Title:    meta.Title,
			Slug:     meta.Slug,
			Date:     meta.Date,
			Summary:  meta.Summary,
			Tags:     meta.Tags,
			Draft:    meta.Draft,
			BodyMD:   md,
			BodyHTML: template.HTML(htmlBody),
		})
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts, nil
}

func splitFrontMatter(raw string) (frontMatter, string, error) {
	var meta frontMatter

	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "---") {
		return meta, "", fmt.Errorf("missing frontmatter delimiter")
	}

	parts := strings.SplitN(trimmed, "---", 3)
	if len(parts) != 3 {
		return meta, "", fmt.Errorf("invalid frontmatter structure")
	}

	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return meta, "", fmt.Errorf("decode yaml frontmatter: %w", err)
	}

	if meta.Title == "" || meta.Slug == "" || meta.Date == "" {
		return meta, "", fmt.Errorf("title, slug, and date are required")
	}

	body := strings.TrimSpace(parts[2])
	return meta, body, nil
}

func renderMarkdown(md string) (string, error) {
	var out bytes.Buffer
	if err := goldmark.Convert([]byte(md), &out); err != nil {
		return "", err
	}
	return out.String(), nil
}
