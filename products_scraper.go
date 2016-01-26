package main

import (
	"math/big"
	"net/url"
)

type Result struct {
	Results     []*Product `json:"results"`
	total       *big.Rat   // otherwise 13.3 + 1.8 = 15.100000000000001
	TotalCached float64    `json:"total"`
}

func newResult() *Result {
	return &Result{
		total: &big.Rat{},
	}
}

type ProductsScraper struct {
	ListScraper *ProductsListScraper
}

func NewProductsScraper(url string) (*ProductsScraper, error) {
	ls, err := NewProductsListScraper(url)
	if err != nil {
		return nil, err
	}

	return &ProductsScraper{
		ListScraper: ls,
	}, nil
}

func (ps *ProductsScraper) Scrape(log AppLogger) (*Result, error) {
	result := newResult()

	links, err := ps.ListScraper.Scrape(log)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		url, err := url.Parse(link)
		if err != nil {
			log.Info(err.Error())
			continue
		}

		f := ps.ListScraper.scraper.f.New(url)
		detailScraper, err := NewProductScraper(f.URL().String())
		if err != nil {
			log.Info(err.Error())
			continue
		}

		product, err := detailScraper.Scrape(log)
		if err != nil {
			log.Info(err.Error())
		}

		result.Results = append(result.Results, product)
		result.total.Add(result.total, product.UnitPrice)
		rounded, _ := result.total.Float64()
		result.TotalCached = rounded
	}

	return result, nil
}
