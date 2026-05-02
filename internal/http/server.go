package httpapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
	"github.com/rs/zerolog"
)

type Server struct {
	port          int
	contentStore  *content.Store
	projectStore  *content.ProjectStore
	aboutPage     content.Page
	siteURL       string
	logger        zerolog.Logger
	searchIndexer *SearchIndexer
}

func New(port int, siteURL string, contentStore *content.Store, projectStore *content.ProjectStore, aboutPage content.Page, logger zerolog.Logger, searchDocs []content.SearchDoc) *http.Server {
	app := &Server{
		port:          port,
		contentStore:  contentStore,
		projectStore:  projectStore,
		aboutPage:     aboutPage,
		siteURL:       siteURL,
		logger:        logger,
		searchIndexer: NewSearchIndexer(searchDocs),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
