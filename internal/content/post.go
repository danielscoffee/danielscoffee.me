package content

import "html/template"

// Post is a published blog article loaded from frontmatter + markdown.
type Post struct {
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Draft   bool

	BodyMD   string
	BodyHTML template.HTML
}
