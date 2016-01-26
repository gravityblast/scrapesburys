package main

type ProductsListScraper struct {
	scraper   *Scraper
	extractor Extractor
}

func NewProductsListScraper(url string) (*ProductsListScraper, error) {
	f, err := NewFetcher(url)
	if err != nil {
		return nil, err
	}

	s := NewScraper(f)
	linkExtractor := NewAttributeExtractor(".product .productInfoWrapper h3 a", "href")
	ext := NewItemExtractor("recipe", DataExtractors{"link": linkExtractor})
	ext.Selector = ".productLister.listView li"

	return &ProductsListScraper{
		scraper:   s,
		extractor: ext,
	}, nil
}

func (ps *ProductsListScraper) Scrape(log AppLogger) ([]string, error) {
	_, items, err := ps.scraper.Scrape(log, Extractors{"products": ps.extractor})
	if err != nil {
		return nil, err
	}

	var links []string

	for _, item := range items {
		links = append(links, item["link"])
	}

	return links, nil
}
