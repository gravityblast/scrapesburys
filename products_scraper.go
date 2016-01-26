package main

import (
	"fmt"
	"math/big"
	"net/url"
	"sync"
)

type Result struct {
	sync.Mutex
	Results     []*Product `json:"results"`
	total       *big.Rat   // otherwise 13.3 + 1.8 = 15.100000000000001
	TotalCached float64    `json:"total"`
}

func (r *Result) Add(product *Product) {
	r.Lock()
	r.Results = append(r.Results, product)
	r.total.Add(r.total, product.UnitPrice)
	rounded, _ := r.total.Float64()
	r.TotalCached = rounded
	r.Unlock()
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

func (ps *ProductsScraper) Scrape(log AppLogger, concurrency int) (*Result, error) {
	result := newResult()

	links, err := ps.ListScraper.Scrape(log)
	if err != nil {
		return nil, err
	}

	c := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			log.Debug(fmt.Sprintf("starting scraper worker %d", i))
			for link := range c {
				product, err := ps.scrapeDetail(log, link)
				if err != nil {
					log.Info(err.Error())
					return
				}
				result.Add(product)
			}
			wg.Done()
		}(i)
	}

	for _, link := range links {
		c <- link
	}
	close(c)

	wg.Wait()

	return result, nil
}

func (ps *ProductsScraper) scrapeDetail(log AppLogger, link string) (*Product, error) {
	url, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	f := ps.ListScraper.scraper.f.New(url)
	detailScraper, err := NewProductScraper(f.URL().String())
	if err != nil {
		return nil, err
	}

	product, err := detailScraper.Scrape(log)
	if err != nil {
		return nil, err
	}

	return product, nil
}
