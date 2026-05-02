package content

import "html/template"

type Project struct {
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Draft   bool

	BodyMD   string
	BodyHTML template.HTML
}
