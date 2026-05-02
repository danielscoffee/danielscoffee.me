package httpapp

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/danielscoffee/danielscoffee.me/internal/web"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/blog", http.StatusFound)
}

func (s *Server) blogIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, web.BlogIndexPage(s.contentStore.All()))
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, web.AboutPage(s.aboutPage))
}

func (s *Server) projectsIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, web.ProjectsIndexPage(s.projectStore.All()))
}

func (s *Server) projectDetailHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/project/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	project, ok := s.projectStore.BySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	s.renderComponent(w, r, web.ProjectDetailPage(project))
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

	s.renderComponent(w, r, web.BlogPostPage(post))
}

func (s *Server) tagIndexHandler(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, "/tag/")
	if tag == "" {
		http.NotFound(w, r)
		return
	}

	posts := s.contentStore.ByTag(tag)
	s.renderComponent(w, r, web.TagPage(tag, posts))
}

func (s *Server) renderComponent(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		s.logger.Error().Err(err).Msg("render component failed")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
