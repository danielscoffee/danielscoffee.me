package content

import "strings"

type sluggable interface {
	slugKey() string
}

type taggable interface {
	tagKeys() []string
}

type searchable interface {
	searchDoc(string) SearchDoc
}

func cloneSlice[T any](items []T) []T {
	return append([]T(nil), items...)
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func buildSlugIndex[T sluggable](items []T) map[string]T {
	index := make(map[string]T, len(items))
	for _, item := range items {
		index[normalizeKey(item.slugKey())] = item
	}
	return index
}

func buildTagIndex[T interface {
	sluggable
	taggable
}](items []T) map[string][]T {
	index := make(map[string][]T)
	for _, item := range items {
		for _, tag := range item.tagKeys() {
			normalized := normalizeKey(tag)
			if normalized == "" {
				continue
			}
			index[normalized] = append(index[normalized], item)
		}
	}
	return index
}

func appendSearchDocs[T searchable](docs []SearchDoc, items []T, docType string) []SearchDoc {
	for _, item := range items {
		docs = append(docs, item.searchDoc(docType))
	}
	return docs
}
