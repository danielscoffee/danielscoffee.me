package content

import (
	"strings"
	"testing"
)

func TestParseNorg_MixedContentAndTasks(t *testing.T) {
	raw := `@document.meta
title: Norg Post
slug: norg-post
date: 2026-05-01
summary: test
tags:
  - norg
  - notes
draft: false
@end
* Heading
*** TODO ship parser
*** DOING write tests
*** DONE pass tests
*** CANCELLED old idea
- bullet
1. ordered
` + "```go\nfmt.Println(\"ok\")\n```" + `
[docs](https://example.com)
`

	meta, body, html, err := parseNorg(raw)
	if err != nil {
		t.Fatalf("parseNorg error: %v", err)
	}
	if meta.Slug != "norg-post" {
		t.Fatalf("slug mismatch: %q", meta.Slug)
	}
	if body == "" || html == "" {
		t.Fatalf("expected non-empty body/html")
	}
	for _, s := range []string{"todo", "doing", "done", "cancelled"} {
		if !strings.Contains(html, `data-task-state="`+s+`"`) {
			t.Fatalf("missing task state %s in html: %s", s, html)
		}
	}
	if !strings.Contains(html, `<span`) {
		t.Fatalf("expected highlighted code spans in html: %s", html)
	}
	if strings.Contains(html, "<html>") || strings.Contains(html, "<body") {
		t.Fatalf("expected code fragment html only, got full document: %s", html)
	}
}

func TestParseNorg_MissingEndFails(t *testing.T) {
	_, _, _, err := parseNorg("@document.meta\ntitle: X\n")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseNorg_UnclosedFenceFails(t *testing.T) {
	raw := `@document.meta
title: X
slug: x
date: 2026-05-01
@end
` + "```go\nfmt.Println(1)"
	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected unclosed code fence error")
	}
}

func TestParseNorg_InvalidTaskStateFails(t *testing.T) {
	raw := `@document.meta
title: X
slug: x
date: 2026-05-01
@end
*** WAITING no state support
`

	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected invalid task state error")
	}
	if !strings.Contains(err.Error(), "invalid task state") {
		t.Fatalf("unexpected error: %v", err)
	}
}
