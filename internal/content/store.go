package content

// Store provides read-only lookups over loaded posts.
type Store struct {
	posts  []Post
	bySlug map[string]Post
	byTag  map[string][]Post
}

func NewStore(posts []Post) *Store {
	postCopies := cloneSlice(posts)
	return &Store{posts: postCopies, bySlug: buildSlugIndex(postCopies), byTag: buildTagIndex(postCopies)}
}

func (s *Store) All() []Post {
	return append([]Post(nil), s.posts...)
}

func (s *Store) Latest(limit int) []Post {
	if limit <= 0 || limit >= len(s.posts) {
		return s.All()
	}
	return append([]Post(nil), s.posts[:limit]...)
}

func (s *Store) BySlug(slug string) (Post, bool) {
	post, ok := s.bySlug[slug]
	return post, ok
}

func (s *Store) ByTag(tag string) []Post {
	normalized := normalizeKey(tag)
	posts := s.byTag[normalized]
	return cloneSlice(posts)
}
