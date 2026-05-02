package content

type SearchDoc struct {
	Type    string
	Title   string
	Slug    string
	Date    string
	Summary string
	Tags    []string
	Body    string
}

func BuildSearchDocs(posts []Post, projects []Project) []SearchDoc {
	docs := make([]SearchDoc, 0, len(posts)+len(projects))
	for _, post := range posts {
		docs = append(docs, SearchDoc{
			Type:    "blog",
			Title:   post.Title,
			Slug:    post.Slug,
			Date:    post.Date,
			Summary: post.Summary,
			Tags:    append([]string(nil), post.Tags...),
			Body:    post.BodyMD,
		})
	}
	for _, project := range projects {
		docs = append(docs, SearchDoc{
			Type:    "projects",
			Title:   project.Title,
			Slug:    project.Slug,
			Date:    project.Date,
			Summary: project.Summary,
			Tags:    append([]string(nil), project.Tags...),
			Body:    project.BodyMD,
		})
	}
	return docs
}
