package httpapp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danielscoffee/danielscoffee.me/internal/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	mux.HandleFunc("/", s.homeHandler)
	mux.HandleFunc("/blog", s.blogIndexHandler)
	mux.HandleFunc("/post/", s.postDetailHandler)
	mux.HandleFunc("/tag/", s.tagIndexHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/rss.xml", s.rssHandler)
	mux.HandleFunc("/sitemap.xml", s.sitemapHandler)
	mux.HandleFunc("/robots.txt", s.robotsHandler)

	return s.securityHeadersMiddleware(mux)
}

func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "up"}
	payload, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(payload); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
