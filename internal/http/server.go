package httpapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

type Server struct {
	port         int
	contentStore *content.Store
	siteURL      string
}

func New(port int, siteURL string, contentStore *content.Store) *http.Server {
	app := &Server{
		port:         port,
		contentStore: contentStore,
		siteURL:      siteURL,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
