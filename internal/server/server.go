package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/danielscoffee/blog/internal/content"
)

type Server struct {
	port         int
	contentStore *content.Store
	siteURL      string
}

func NewServer() *http.Server {
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
		log.Fatalf("load posts: %v", err)
	}

	app := &Server{
		port:         port,
		contentStore: content.NewStore(posts),
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
