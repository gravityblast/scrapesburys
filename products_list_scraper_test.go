package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductsListScraper(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	log := newAppLogger(logLevelFatal)
	pls, err := NewProductsListScraper("/index.html")
	require.Nil(err)

	pls.scraper.f = newTestHTTPFetcher("index.html")
	links, err := pls.Scrape(log)
	require.Nil(err)

	expected := []string{"1.html", "2.html"}
	assert.Equal(expected, links)
}
