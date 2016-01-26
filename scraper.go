package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	f Fetcher
}

func NewScraper(f Fetcher) *Scraper {
	return &Scraper{
		f: f,
	}
}

func (s *Scraper) Scrape(log AppLogger, exts Extractors) (int64, Items, error) {
	var items Items

	resp := s.f.Fetch(log)
	if resp.Error() != nil {
		return 0, nil, resp.Error()
	}
	r := resp.ReadCloser()
	defer r.Close()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return 0, nil, err
	}

	for name, ext := range exts {
		log.Debug(fmt.Sprintf("extract `%s` from `%s`", name, s.f.URL()))
		items = append(items, ext.Extract(doc.Selection)...)
	}

	return resp.Length(), items, nil
}
