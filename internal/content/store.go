package content

import "strings"

// Store provides read-only lookups over loaded posts.
type Store struct {
	posts  []Post
	bySlug map[string]Post
	byTag  map[string][]Post
}

func NewStore(posts []Post) *Store {
	postCopies := append([]Post(nil), posts...)
	bySlug := make(map[string]Post, len(postCopies))
	byTag := make(map[string][]Post)

	for _, post := range postCopies {
		bySlug[post.Slug] = post
		for _, tag := range post.Tags {
			normalized := strings.ToLower(strings.TrimSpace(tag))
			if normalized == "" {
				continue
			}
			byTag[normalized] = append(byTag[normalized], post)
		}
	}

	return &Store{posts: postCopies, bySlug: bySlug, byTag: byTag}
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
	normalized := strings.ToLower(strings.TrimSpace(tag))
	posts := s.byTag[normalized]
	return append([]Post(nil), posts...)
}
