package content

import "html/template"

type Page struct {
	Title   string
	Slug    string
	Date    string
	Summary string

	BodyMD   string
	BodyHTML template.HTML
}
