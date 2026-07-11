package httpapp

import (
	"compress/gzip"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", staticCacheMiddleware(fileServer))

	mux.HandleFunc("/", s.homeHandler)
	mux.HandleFunc("/blog", s.blogIndexHandler)
	mux.HandleFunc("/about", s.aboutHandler)
	mux.HandleFunc("/projects", s.projectsIndexHandler)
	mux.HandleFunc("/projects/", s.projectsTreeHandler)
	mux.HandleFunc("/project/", s.legacyProjectRedirectHandler)
	mux.HandleFunc("/post/", s.postDetailHandler)
	mux.HandleFunc("/tag/", s.tagIndexHandler)
	mux.HandleFunc("/search", s.searchHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/rss.xml", s.rssHandler)
	mux.HandleFunc("/sitemap.xml", s.sitemapHandler)
	mux.HandleFunc("/robots.txt", s.robotsHandler)

	return s.requestLoggingMiddleware(s.securityHeadersMiddleware(s.compressionMiddleware(mux)))
}

func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Permissions-Policy", "camera=(), geolocation=(), microphone=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' https: data:; font-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'")
		next.ServeHTTP(w, r)
	})
}

func staticCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) compressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.Header.Get("Range") != "" {
			next.ServeHTTP(w, r)
			return
		}

		appendVary(w.Header(), "Accept-Encoding")
		gz := gzip.NewWriter(w)
		grw := &gzipResponseWriter{ResponseWriter: w, writer: gz}
		next.ServeHTTP(grw, r)
		if grw.Header().Get("Content-Encoding") == "gzip" {
			if err := gz.Close(); err != nil {
				s.logger.Error().Err(err).Msg("close gzip response failed")
			}
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer      *gzip.Writer
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	if shouldCompressStatus(code) {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.Header().Get("Content-Encoding") != "gzip" {
		return w.ResponseWriter.Write(b)
	}
	return w.writer.Write(b)
}

func shouldCompressStatus(code int) bool {
	return code >= 200 && code != http.StatusNoContent && code != http.StatusNotModified
}

func appendVary(h http.Header, value string) {
	current := h.Values("Vary")
	for _, existing := range current {
		for part := range strings.SplitSeq(existing, ",") {
			if strings.EqualFold(strings.TrimSpace(part), value) {
				return
			}
		}
	}
	h.Add("Vary", value)
}

func (s *Server) requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		s.logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lrw.statusCode).
			Int("bytes", lrw.bytes).
			Dur("duration", time.Since(start)).
			Str("remote_ip", clientIP(r)).
			Str("user_agent", r.UserAgent()).
			Msg("http_request")
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "up"}
	payload, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error().Err(err).Msg("marshal health response failed")
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(payload); err != nil {
		s.logger.Error().Err(err).Msg("write health response failed")
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

func clientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
