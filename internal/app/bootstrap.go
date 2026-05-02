package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
	httpapp "github.com/danielscoffee/danielscoffee.me/internal/http"
)

func NewServer() (*http.Server, error) {
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
		return nil, fmt.Errorf("load posts: %w", err)
	}

	store := content.NewStore(posts)
	return httpapp.New(port, siteURL, store), nil
}
