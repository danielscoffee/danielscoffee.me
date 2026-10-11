package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

type postYear struct {
	Year  string
	Posts []content.Post
}

type projectYear struct {
	Year     string
	Projects []content.Project
}

// postsByYear groups posts (already newest first) under their year.
func postsByYear(posts []content.Post) []postYear {
	var groups []postYear
	for _, post := range posts {
		year := yearOf(post.Date)
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, postYear{Year: year})
		}
		groups[len(groups)-1].Posts = append(groups[len(groups)-1].Posts, post)
	}
	return groups
}

// projectsByYear groups projects (already newest first) under their year.
func projectsByYear(projects []content.Project) []projectYear {
	var groups []projectYear
	for _, project := range projects {
		year := yearOf(project.Date)
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, projectYear{Year: year})
		}
		groups[len(groups)-1].Projects = append(groups[len(groups)-1].Projects, project)
	}
	return groups
}

func yearOf(date string) string {
	if len(date) >= 4 {
		return date[:4]
	}
	return date
}

// shortDate turns 2026-10-10 into "Oct 10"; the year heading carries the year.
func shortDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("Jan 2")
}

// longDate turns 2026-10-10 into "Oct 10, 2026" for lists without year marks.
func longDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("Jan 2, 2006")
}

// readingTime estimates minutes at 200 words a minute from the source text.
func readingTime(body string) string {
	minutes := max(1, (len(strings.Fields(body))+199)/200)
	return fmt.Sprintf("%d min", minutes)
}
