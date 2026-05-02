# danielscoffee.me

Personal site + blog built with Go, templ, and Tailwind.

## Stack

- Go (`net/http`)
- [templ](https://templ.guide/) for server-rendered pages
- Tailwind CSS
- Markdown posts with YAML frontmatter (stored in Git)

## Local development

1. Copy env file:

```bash
cp .env.example .env
```

2. Build and run:

```bash
make build
make run
```

The app runs on `http://localhost:8080` by default.

## Writing posts

Add Markdown files in `content/posts/*.md`:

```md
---
title: "My Post"
slug: "my-post"
date: "2026-04-26"
summary: "Short summary"
tags: ["go", "personal"]
draft: false
---
# Post title

Body content...
```

Push commits to publish updates.

## Useful commands

```bash
make generate   # templ + tailwind
make test       # generate + go test ./...
make docker-run # run with Docker Compose
make docker-down
```

## Routes

- `/` home
- `/blog` all posts
- `/post/{slug}` post detail
- `/tag/{tag}` posts by tag
- `/rss.xml`
- `/sitemap.xml`
- `/robots.txt`
- `/health`
