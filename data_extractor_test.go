package main

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func newTestSelection(t *testing.T, html string) *goquery.Selection {
	r := bytes.NewBufferString(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatal(err)
	}

	return doc.Selection
}

func TestTextExtractor_Extract(t *testing.T) {
	sel := newTestSelection(t, "<html><body><h2 class=\"name\">\n\tItem 1\r\n</h2></body></html>")
	e := NewTextExtractor(".name")
	assert.Equal(t, "Item 1", e.Extract(sel))
}

func TestTextExtractor_Extract_NTH(t *testing.T) {
	sel := newTestSelection(t, `
		<ul>
			<li>One</li>
			<li>Two</li>
		</ul>
	`)
	e := NewTextExtractor("ul li:first-child")
	assert.Equal(t, "One", e.Extract(sel))

	sel = newTestSelection(t, `
	<HTMLContent>
		<h3 class="productDataItemHeader">Description 1</h3>
		<div class="productText">
			<p>Content 1</p>
			<p><p></p></p>
		</div>
		<h3 class="productDataItemHeader">Description 2</h3>
		<div class="productText">
			<p>Content 2</p>
			<p><p></p></p>
		</div>
	</HTMLContent>
	`)
	e = NewTextExtractor("HTMLContent div:first-of-type p")
	assert.Equal(t, "Content 1", e.Extract(sel))
}

func TestAttributeExtractor_Extract(t *testing.T) {
	sel := newTestSelection(t, `<hmtl><body><a href="/foo">Hello</a></body></html>`)
	e := NewAttributeExtractor("a", "href")
	assert.Equal(t, "/foo", e.Extract(sel))
}
