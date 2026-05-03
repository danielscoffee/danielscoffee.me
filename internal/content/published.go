package content

type Published struct {
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Draft   bool
}

func (p Published) publishedDate() string {
	return p.Date
}

func (p Published) slugKey() string {
	return p.Slug
}

func (p Published) tagKeys() []string {
	return p.Tags
}
