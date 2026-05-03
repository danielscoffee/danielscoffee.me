package content

import "html/template"

// Post is a published blog article loaded from frontmatter + markdown.
type Post struct {
	Published

	BodyMD   string
	BodyHTML template.HTML
}

func (p Post) searchDoc(docType string) SearchDoc {
	return SearchDoc{
		Type:    docType,
		Title:   p.Title,
		Slug:    p.Slug,
		Date:    p.Date,
		Summary: p.Summary,
		Tags:    append([]string(nil), p.Tags...),
		Body:    p.BodyMD,
	}
}
