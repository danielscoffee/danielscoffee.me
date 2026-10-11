package httpapp

import (
	"net/http"
	"strings"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

// llmsHandler serves /llms.txt (https://llmstxt.org): a Markdown map of the
// site for language models and agents, built from the same stores as the
// sitemap so it never drifts from the published content.
func (s *Server) llmsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(buildLLMSTxt(s.siteURL, s.contentStore.All(), s.projectStore.All(), s.aboutPage)))
}

func buildLLMSTxt(siteURL string, posts []content.Post, projects []content.Project, about content.Page) string {
	base := strings.TrimRight(siteURL, "/")
	var b strings.Builder

	b.WriteString("# danielscoffee\n\n")
	b.WriteString("> Personal site and blog engine of Daniel, a backend developer: posts on Go and backend systems, and devlogs for the Go, templ, and Neorg engine that runs this site.\n\n")
	b.WriteString("Daniel works on production healthcare software (Java, PL/SQL, Kafka, AWS, Docker, event-driven systems) and builds Go infrastructure, devtools, and distributed-systems experiments outside work. Posts and devlogs are written in Neorg and served as server-rendered HTML.\n")

	if len(posts) > 0 {
		b.WriteString("\n## Blog\n\n")
		for _, post := range posts {
			writeLLMSLink(&b, post.Title, base+"/post/"+post.Slug, post.Summary)
		}
	}

	if len(projects) > 0 {
		b.WriteString("\n## Projects\n\n")
		for _, project := range projects {
			projectURL := base + "/projects/" + project.Slug
			writeLLMSLink(&b, project.Title, projectURL, project.Summary)
			for _, sub := range project.SubPosts {
				writeLLMSLink(&b, sub.Title, projectURL+"/"+sub.Slug, devlogNote(project.Title, sub.Summary))
			}
		}
	}

	if about.Title != "" {
		b.WriteString("\n## About\n\n")
		writeLLMSLink(&b, about.Title, base+"/about", about.Summary)
	}

	b.WriteString("\n## Optional\n\n")
	writeLLMSLink(&b, "RSS feed", base+"/rss.xml", "Posts, projects, and devlogs, newest first.")
	writeLLMSLink(&b, "Sitemap", base+"/sitemap.xml", "Every public URL.")
	return b.String()
}

func devlogNote(project, summary string) string {
	if summary == "" {
		return project + " devlog."
	}
	return project + " devlog: " + summary
}

// writeLLMSLink writes one "- [title](url): note" line, escaping the
// brackets that would break the Markdown link.
func writeLLMSLink(b *strings.Builder, title, url, note string) {
	title = strings.NewReplacer(`\`, `\\`, "[", `\[`, "]", `\]`).Replace(strings.TrimSpace(title))
	b.WriteString("- [")
	b.WriteString(title)
	b.WriteString("](")
	b.WriteString(url)
	b.WriteString(")")
	if note = strings.TrimSpace(note); note != "" {
		b.WriteString(": ")
		b.WriteString(strings.ReplaceAll(note, "\n", " "))
	}
	b.WriteString("\n")
}
