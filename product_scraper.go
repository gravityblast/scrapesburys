package main

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
)

type Product struct {
	Title           string   `json:"title"`
	Size            string   `json:"size"`
	UnitPrice       *big.Rat `json:"-"`
	UnitPriceCached float64  `json:"unit_price"`
	Description     string   `json:"description"`
}

type ProductScraper struct {
	scraper   *Scraper
	extractor Extractor
}

func NewProductScraper(url string) (*ProductScraper, error) {
	f, err := NewFetcher(url)
	if err != nil {
		return nil, err
	}

	s := NewScraper(f)
	title := NewTextExtractor(".productSummary .productTitleDescriptionContainer h1")
	unitPrice := NewTextExtractor(".productSummary .addToTrolleytabBox .pricePerUnit")
	description := NewTextExtractor(".mainProductInfoWrapper .mainProductInfo .tabs #information productcontent div.productText:first-of-type")

	ext := NewItemExtractor("product", DataExtractors{"title": title, "unit_price": unitPrice, "description": description})
	ext.Selector = ".section.productContent"

	return &ProductScraper{
		scraper:   s,
		extractor: ext,
	}, nil
}

func (ps *ProductScraper) Scrape(log AppLogger) (*Product, error) {
	length, items, err := ps.scraper.Scrape(log, Extractors{"products": ps.extractor})
	if err != nil {
		return nil, err
	}

	if len(items) < 1 {
		return nil, errors.New("no data found")
	}

	item := items[0]

	price, priceFloat := parsePriceRat(item["unit_price"])

	p := &Product{
		Title:           item["title"],
		UnitPrice:       price,
		UnitPriceCached: priceFloat,
		Description:     item["description"],
		Size:            fmt.Sprintf("%fkb", float64(length)/1024.0),
	}

	return p, nil
}

var priceRegexp = regexp.MustCompile(`[^0-9]*([0-9\.]+)`)

func parsePriceRat(s string) (*big.Rat, float64) {
	var number string
	matches := priceRegexp.FindStringSubmatch(s)
	if len(matches) == 2 {
		number = matches[1]
	}

	f := big.NewFloat(0)
	f.Parse(number, 10)
	r, _ := f.Rat(nil)
	rounded, _ := r.Float64()

	return r, rounded
}
