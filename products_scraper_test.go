package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductsScraper(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	ts := httptest.NewServer(http.FileServer(http.Dir("fixtures")))

	defer ts.Close()

	log := newAppLogger(logLevelFatal)
	ps, err := NewProductsScraper(fmt.Sprintf("%s/index.html", ts.URL))
	require.Nil(err)
	result, err := ps.Scrape(log, 1)
	require.Nil(err)

	assert.Len(result.Results, 2)
	assert.Equal(5.0, result.TotalCached)
}
