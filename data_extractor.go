package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DataExtractor is an interface that wraps simple exatractors behaviours used to extract text
type DataExtractor interface {
	Extract(*goquery.Selection) string
}

// DataExtractors is a collection of named DataExtractors
type DataExtractors map[string]DataExtractor

// TextExtractor is an extractor used to extract text from a selection
type TextExtractor struct {
	Selector string
}

// NewTextExtractor returns a TextExtractor with a given selector
func NewTextExtractor(selector string) *TextExtractor {
	return &TextExtractor{
		Selector: selector,
	}
}

// Extract extracts text of elements returned querying with TextExtractor's selector
func (e *TextExtractor) Extract(s *goquery.Selection) string {
	text := s.Find(e.Selector).Text()
	return strings.TrimSpace(text)
}

// AttributeExtractor is an exatractor used to extract values from an element's attribute
type AttributeExtractor struct {
	Selector      string
	AttributeName string
}

// NewAttributeExtractor returns a new AttributeExtractor
func NewAttributeExtractor(selector string, attributeName string) *AttributeExtractor {
	return &AttributeExtractor{
		Selector:      selector,
		AttributeName: attributeName,
	}
}

// Extract extracts the value of a specific attribute
func (e *AttributeExtractor) Extract(s *goquery.Selection) string {
	v, _ := s.Find(e.Selector).Attr(e.AttributeName)
	return v
}
