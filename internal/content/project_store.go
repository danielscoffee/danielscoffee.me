package content

import "strings"

type ProjectStore struct {
	projects []Project
	bySlug   map[string]Project
}

func NewProjectStore(projects []Project) *ProjectStore {
	copies := append([]Project(nil), projects...)
	bySlug := make(map[string]Project, len(copies))
	for _, project := range copies {
		bySlug[strings.ToLower(strings.TrimSpace(project.Slug))] = project
	}
	return &ProjectStore{projects: copies, bySlug: bySlug}
}

func (s *ProjectStore) All() []Project {
	return append([]Project(nil), s.projects...)
}

func (s *ProjectStore) BySlug(slug string) (Project, bool) {
	project, ok := s.bySlug[strings.ToLower(strings.TrimSpace(slug))]
	return project, ok
}
