package httpapp

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/danielscoffee/danielscoffee.me/internal/content"
)

func (s *Server) rssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(buildRSS(s.siteURL, s.contentStore.All())))
}

func (s *Server) sitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(buildSitemap(s.siteURL, s.contentStore.All())))
}

func (s *Server) robotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintf(w, "User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", strings.TrimRight(s.siteURL, "/"))
}

type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description,omitempty"`
}

func buildRSS(siteURL string, posts []content.Post) string {
	base := strings.TrimRight(siteURL, "/")
	items := make([]rssItem, 0, len(posts))
	for _, post := range posts {
		items = append(items, rssItem{
			Title:       post.Title,
			Link:        base + path.Join("/post", post.Slug),
			PubDate:     toRFC1123(post.Date),
			Description: post.Summary,
		})
	}

	doc := rssDocument{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "Daniel's Site",
			Link:        base,
			Description: "Personal blog",
			Items:       items,
		},
	}

	payload, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"></rss>`
	}
	return xml.Header + string(payload)
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

func buildSitemap(siteURL string, posts []content.Post) string {
	base := strings.TrimRight(siteURL, "/")

	urls := []sitemapURL{{Loc: base + "/"}, {Loc: base + "/blog"}, {Loc: base + "/rss.xml"}}
	for _, post := range posts {
		urls = append(urls, sitemapURL{Loc: base + path.Join("/post", post.Slug), LastMod: post.Date})
	}

	payload, err := xml.MarshalIndent(sitemapURLSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls}, "", "  ")
	if err != nil {
		return `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`
	}
	return xml.Header + string(payload)
}

func toRFC1123(isoDate string) string {
	t, err := time.Parse("2006-01-02", isoDate)
	if err != nil {
		return isoDate
	}
	return t.Format(time.RFC1123Z)
}
