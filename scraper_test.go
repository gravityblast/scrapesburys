package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScraper(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	log := newAppLogger(LogLevelFatal)
	f := newTestHTTPFetcher("/hello.html")
	s := NewScraper(f)

	linkExtractor := NewAttributeExtractor("a", "href")
	itemsExtractor := NewItemExtractor("recipe", DataExtractors{"link": linkExtractor})
	itemsExtractor.Selector = ".list li"

	_, items, err := s.Scrape(log, Extractors{"items": itemsExtractor})
	require.Nil(err)

	expected := Items{
		{"link": "foo.html"},
		{"link": "bar.html"},
		{"link": "baz.html"},
	}
	assert.Equal(expected, items)
}
