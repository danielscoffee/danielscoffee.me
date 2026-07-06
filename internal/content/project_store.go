package content

type ProjectStore struct {
	projects []Project
	bySlug   map[string]Project
	byTag    map[string][]Project
}

func NewProjectStore(projects []Project) *ProjectStore {
	copies := cloneSlice(projects)
	return &ProjectStore{projects: copies, bySlug: buildSlugIndex(copies), byTag: buildTagIndex(copies)}
}

func (s *ProjectStore) All() []Project {
	return append([]Project(nil), s.projects...)
}

func (s *ProjectStore) BySlug(slug string) (Project, bool) {
	project, ok := s.bySlug[normalizeKey(slug)]
	return project, ok
}

func (s *ProjectStore) ByTag(tag string) []Project {
	projects := s.byTag[normalizeKey(tag)]
	return cloneSlice(projects)
}

// SubPost returns the named subpost within a project. Both lookups are case-insensitive.
func (s *ProjectStore) SubPost(projectSlug, subSlug string) (Project, ProjectSubPost, bool) {
	project, ok := s.BySlug(projectSlug)
	if !ok {
		return Project{}, ProjectSubPost{}, false
	}
	normalized := normalizeKey(subSlug)
	for _, sub := range project.SubPosts {
		if normalizeKey(sub.Slug) == normalized {
			return project, sub, true
		}
	}
	return Project{}, ProjectSubPost{}, false
}
