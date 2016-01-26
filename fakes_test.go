package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type fakeClient struct {
	httpClient *httpClient
}

func (fc *fakeClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{}
	root := "./fixtures"
	f, err := os.Open(filepath.Join(root, filepath.FromSlash(path.Clean(req.URL.Path))))
	if err != nil {
		return nil, err
	}

	resp.Body = f

	return resp, nil
}

func newTestHTTPFetcher(path string) Fetcher {
	url, _ := url.Parse(fmt.Sprintf("http://localhost:0000/%s", path))
	f := NewHTTPFetcher(url)
	f.client = &fakeClient{}

	return f
}
