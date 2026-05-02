package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
	httpapp "github.com/danielscoffee/danielscoffee.me/internal/http"
	"github.com/danielscoffee/danielscoffee.me/internal/logging"
)

type Runtime struct {
	Server       *http.Server
	Logger       zerolog.Logger
	Port         int
	SiteURL      string
	PostCount    int
	ProjectCount int
	LogCfg       logging.Config
}

func NewRuntime() (*Runtime, error) {
	logCfg := logging.ConfigFromEnv()
	logger := logging.New(logCfg, os.Stdout)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	siteURL := strings.TrimRight(os.Getenv("SITE_URL"), "/")
	if siteURL == "" {
		siteURL = fmt.Sprintf("http://localhost:%d", port)
	}

	posts, err := content.LoadPosts("content/posts")
	if err != nil {
		logger.Error().Err(err).Msg("load posts failed")
		return nil, fmt.Errorf("load posts: %w", err)
	}

	projects, err := content.LoadProjects("content/projects")
	if err != nil {
		logger.Error().Err(err).Msg("load projects failed")
		return nil, fmt.Errorf("load projects: %w", err)
	}

	aboutPage, err := content.LoadPage("content/pages/about.norg")
	if err != nil {
		logger.Error().Err(err).Msg("load about page failed")
		return nil, fmt.Errorf("load about page: %w", err)
	}

	logger.Info().
		Int("post_count", len(posts)).
		Int("project_count", len(projects)).
		Str("content_dir", "content/").
		Msg("content loaded")

	postStore := content.NewStore(posts)
	projectStore := content.NewProjectStore(projects)
	searchDocs := content.BuildSearchDocs(posts, projects)

	apiServer := httpapp.New(port, siteURL, postStore, projectStore, aboutPage, logger, searchDocs)

	return &Runtime{
		Server:       apiServer,
		Logger:       logger,
		Port:         port,
		SiteURL:      siteURL,
		PostCount:    len(posts),
		ProjectCount: len(projects),
		LogCfg:       logCfg,
	}, nil
}
