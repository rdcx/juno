package client

import (
	"bytes"
	"encoding/json"
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/crawl/dto"
	"net/http"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) Crawl(url string) error {

	crawlReq := dto.CrawlRequest{URL: url}

	jsonB, err := json.Marshal(crawlReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/crawl", bytes.NewBuffer(jsonB))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return crawl.ErrFailedCrawlRequest
	}

	return nil
}

func (c *Client) CrawlURLs(urls []string) error {

	crawlReq := dto.CrawlURLsRequest{URLs: urls}

	jsonB, err := json.Marshal(crawlReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/crawl/urls", bytes.NewBuffer(jsonB))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return crawl.ErrFailedCrawlRequest
	}

	return nil
}
