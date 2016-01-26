package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductScraper(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	log := newAppLogger(LogLevelFatal)
	ps, err := NewProductScraper("/1.html")
	require.Nil(err)

	ps.scraper.f = newTestHTTPFetcher("1.html")
	p, err := ps.Scrape(log)
	require.Nil(err)

	assert.Equal("Sainsbury's Apricot Ripe & Ready x5", p.Title)
	assert.Equal("0.000000kb", p.Size)
	assert.Equal(3.5, p.UnitPriceCached)
	assert.Equal("Apricots", p.Description)
}
