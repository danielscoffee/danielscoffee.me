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
@table
| Name | Type |
| - | - |
| Search | blog |
| Projects | projects |
@end
#priority high
> this is quote line
$ term: definition body
*bold* /italic/ _underline_ !spoiler! $x^2 + y^2$
@code go
fmt.Println("ok")
@end
[docs](https://example.com)
![diagram](https://cdn.example.com/images/diagram.png)
.image https://cdn.example.com/images/screenshot.png
{https://cdn.example.com/images/asset.png}[asset image](image)
!{https://cdn.example.com/images/old.png}[old image]
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
	if !strings.Contains(html, `<pre class="chroma" tabindex="0">`) {
		t.Fatalf("expected keyboard-focusable highlighted code, got %s", html)
	}
	for _, marker := range []string{
		`<blockquote data-priority="high"><p>this is quote line</p></blockquote>`,
		`<dl class="norg-definitions">`,
		`<dt>term</dt>`,
		`<dd>definition body</dd>`,
		`<strong>bold</strong>`,
		`<em>italic</em>`,
		`<u>underline</u>`,
		`<span class="spoiler">spoiler</span>`,
		`<span class="math-latex">x^2 + y^2</span>`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("expected rendered marker %q, got %s", marker, html)
		}
	}
	if !strings.Contains(html, `<table tabindex="0">`) || !strings.Contains(html, `<th>Name</th>`) || !strings.Contains(html, `<td>blog</td>`) {
		t.Fatalf("expected keyboard-focusable table html, got %s", html)
	}
	if strings.Contains(html, "<html>") || strings.Contains(html, "<body") {
		t.Fatalf("expected code fragment html only, got full document: %s", html)
	}
	for _, marker := range []string{
		`<img class="post-image" src="https://cdn.example.com/images/diagram.png" alt="diagram"`,
		`<img class="post-image" src="https://cdn.example.com/images/screenshot.png" alt="screenshot"`,
		`<img class="post-image" src="https://cdn.example.com/images/asset.png" alt="asset image"`,
		`<img class="post-image" src="https://cdn.example.com/images/old.png" alt="old image"`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("expected rendered image marker %q, got %s", marker, html)
		}
	}
}

func TestParseNorg_OffsetsHeadingsToReservePageTitle(t *testing.T) {
	_, _, html, err := parseNorg(norgWithBody("* Section\n** Subsection"))
	if err != nil {
		t.Fatalf("parseNorg error: %v", err)
	}
	capped, err := renderNorgHTML([]norgNode{
		{kind: norgHeading, level: 5, text: "Deep"},
		{kind: norgHeading, level: 6, text: "Deeper"},
	})
	if err != nil {
		t.Fatalf("renderNorgHTML error: %v", err)
	}
	html += capped
	for _, heading := range []string{
		"<h2>Section</h2>",
		"<h3>Subsection</h3>",
		"<h6>Deep</h6>",
		"<h6>Deeper</h6>",
	} {
		if !strings.Contains(html, heading) {
			t.Errorf("expected %q in %s", heading, html)
		}
	}
	if strings.Contains(html, "<h1") {
		t.Errorf("Norg body must reserve h1 for page title: %s", html)
	}
}

func TestPlainCodeHTML_IsKeyboardFocusable(t *testing.T) {
	html := plainCodeHTML("text", "a < b")
	if !strings.Contains(html, `<pre tabindex="0"><code class="language-text">a &lt; b</code></pre>`) {
		t.Fatalf("expected keyboard-focusable escaped code, got %s", html)
	}
}

func TestParseNorg_PreservesUTF8Punctuation(t *testing.T) {
	raw := "@document.meta\n" +
		"title: UTF8\n" +
		"slug: utf8\n" +
		"date: 2026-05-01\n" +
		"@end\n" +
		"Search is one of the few things on this site that runs on every request. No DB, no Elastic — just a slice of `SearchDoc` and a scorer that walks it.\n"

	_, _, html, err := parseNorg(raw)
	if err != nil {
		t.Fatalf("parseNorg error: %v", err)
	}
	if !strings.Contains(html, "Elastic — just") {
		t.Fatalf("expected em dash preserved, got %s", html)
	}
	if strings.Contains(html, "â") {
		t.Fatalf("expected no mojibake, got %s", html)
	}
}

func TestParseNorg_MissingEndFails(t *testing.T) {
	_, _, _, err := parseNorg("@document.meta\ntitle: X\n")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseNorg_UnclosedCodeBlockFails(t *testing.T) {
	raw := `@document.meta
title: X
slug: x
date: 2026-05-01
@end
@code go
fmt.Println(1)
`
	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected unclosed @code block error")
	}
}

func TestParseNorg_UnclosedTableBlockFails(t *testing.T) {
	raw := `@document.meta
title: X
slug: x
date: 2026-05-01
@end
@table
| A | B |
| - | - |
| 1 | 2 |
`
	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected unclosed @table block error")
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

func TestParseNorg_ValidLinks(t *testing.T) {
	for _, href := range []string{
		"/about",
		"./guide",
		"../index",
		"#section",
		"?page=2",
		"https://example.com/docs",
		"http://example.com/docs",
		"mailto:hello@example.com",
	} {
		t.Run(href, func(t *testing.T) {
			_, _, rendered, err := parseNorg(norgWithBody("[link](" + href + ")"))
			if err != nil {
				t.Fatalf("parse valid link: %v", err)
			}
			if !strings.Contains(rendered, `href="`+href+`"`) {
				t.Fatalf("expected href %q, got %s", href, rendered)
			}
		})
	}
}

func TestParseNorg_InvalidLinks(t *testing.T) {
	for _, href := range []string{
		"javascript:alert(1)",
		"data:text/html,pwn",
		"ftp://example.com/file",
		"//example.com/path",
		"https:///missing-host",
		"bad\x00path",
		"\t/about",
		"/about\t",
		`/\evil.com`,
		`/\\evil.com`,
	} {
		t.Run(href, func(t *testing.T) {
			_, _, _, err := parseNorg(norgWithBody("[link](" + href + ")"))
			if err == nil {
				t.Fatal("expected invalid link error")
			}
			if !strings.Contains(err.Error(), "invalid link url") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func norgWithBody(body string) string {
	return "@document.meta\ntitle: X\nslug: x\ndate: 2026-05-01\n@end\n" + body + "\n"
}

func TestParseNorg_InvalidCDNImageFails(t *testing.T) {
	raw := `@document.meta
title: X
slug: x
date: 2026-05-01
@end
![bad](https://evil.example.com/pwn.png)
`

	_, _, _, err := parseNorg(raw)
	if err == nil {
		t.Fatal("expected CDN image validation error")
	}
	if !strings.Contains(err.Error(), "image host must be CDN") {
		t.Fatalf("unexpected error: %v", err)
	}
}
