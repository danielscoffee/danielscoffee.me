package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/danielscoffee/blog/internal/web"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderComponent(w, r, web.HomePage(s.contentStore.Latest(5)))
}

func (s *Server) blogIndexHandler(w http.ResponseWriter, r *http.Request) {
	renderComponent(w, r, web.BlogIndexPage(s.contentStore.All()))
}

func (s *Server) postDetailHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	post, ok := s.contentStore.BySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	renderComponent(w, r, web.BlogPostPage(post))
}

func (s *Server) tagIndexHandler(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, "/tag/")
	if tag == "" {
		http.NotFound(w, r)
		return
	}

	posts := s.contentStore.ByTag(tag)
	renderComponent(w, r, web.TagPage(tag, posts))
}

func renderComponent(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("render component: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
