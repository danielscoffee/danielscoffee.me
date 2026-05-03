package content

type ProjectStore struct {
	projects []Project
	bySlug   map[string]Project
}

func NewProjectStore(projects []Project) *ProjectStore {
	copies := cloneSlice(projects)
	return &ProjectStore{projects: copies, bySlug: buildSlugIndex(copies)}
}

func (s *ProjectStore) All() []Project {
	return append([]Project(nil), s.projects...)
}

func (s *ProjectStore) BySlug(slug string) (Project, bool) {
	project, ok := s.bySlug[normalizeKey(slug)]
	return project, ok
}
