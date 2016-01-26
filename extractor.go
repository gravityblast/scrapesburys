package main

import "github.com/PuerkitoBio/goquery"

// Extractor is an interface that wraps methods for elements data extractors
type Extractor interface {
	Name() string
	Extract(sel *goquery.Selection) Items
}

// Extractors is a collection of named extractors
type Extractors map[string]Extractor

// Item is a map returned from extractors
type Item map[string]string

// Items is a collection of items
type Items []Item

// ItemExtractor executes multiple DataExtractors and it's used to extract Items
type ItemExtractor struct {
	name           string
	Selector       string
	DataExtractors DataExtractors
}

// NewItemExtractor returns a new ItemExtractor with a given name and extractors
func NewItemExtractor(name string, de DataExtractors) *ItemExtractor {
	return &ItemExtractor{
		name:           name,
		DataExtractors: de,
	}
}

// Extract returns Items querying `sel` using ItemExtractor's DataExtractors
func (e *ItemExtractor) Extract(sel *goquery.Selection) Items {
	var items Items

	if e.Selector != "" {
		sel = sel.Find(e.Selector)
	}

	sel.Each(func(i int, s *goquery.Selection) {
		item := make(Item)
		for name, ext := range e.DataExtractors {
			item[name] = ext.Extract(s)
		}
		items = append(items, item)
	})

	return items

}

// Name returns the ItemExtractor's name
func (e *ItemExtractor) Name() string {
	return e.name
}
