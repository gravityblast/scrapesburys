package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Response wraps an http response adding errors and size of the returned document
type Response interface {
	ReadCloser() io.ReadCloser
	Error() error
	Length() int64
}

type response struct {
	r      io.ReadCloser
	err    error
	length int64
}

// ReadCloser returns the response ReadCloser
func (r *response) ReadCloser() io.ReadCloser {
	return r.r
}

// Error returns errors returned from the http client
func (r *response) Error() error {
	return r.err
}

// Length returns the response size
func (r *response) Length() int64 {
	return r.length
}

// Fetcher is an interface the wraps methods used to retrieve a local or remote document
type Fetcher interface {
	Fetch(AppLogger) Response
	New(*url.URL) Fetcher
	URL() *url.URL
}

// NewFetcher returns a new Fetcher
func NewFetcher(u string) (Fetcher, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return NewHTTPFetcher(url), nil
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPFetcher is a Fetcher specific to retrieve documents via HTTP
type HTTPFetcher struct {
	client httpClient
	u      *url.URL
}

// NewHTTPFetcher returns a new HTTPFetcher
func NewHTTPFetcher(u *url.URL) *HTTPFetcher {
	client := &http.Client{}
	return &HTTPFetcher{
		client: client,
		u:      u,
	}
}

// Fetch retrieves the document and returns a Response object
func (f *HTTPFetcher) Fetch(log AppLogger) Response {
	r := &response{}
	log.Info(fmt.Sprintf("fetching `%s`", f.u.String()))
	req, err := http.NewRequest("GET", f.u.String(), nil)
	if err != nil {
		r.err = err
		log.Error(fmt.Sprintf("%v", err))
		return r
	}

	resp, err := f.client.Do(req)
	if err != nil {
		r.err = err
		log.Error(fmt.Sprintf("%v", err))
		return r
	}

	r.r = resp.Body
	r.length = resp.ContentLength

	return r
}

// URL returns the Fetcher's URL
func (f *HTTPFetcher) URL() *url.URL {
	return f.u
}

// New returns a new Fetcher for a new URL resolving references with the current one
func (f *HTTPFetcher) New(u *url.URL) Fetcher {
	return NewHTTPFetcher(f.u.ResolveReference(u))
}
