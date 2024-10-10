package domain

import "errors"

var ErrFailedCrawlRequest = errors.New("failed to send crawl request")

type CrawlRequest struct {
	URL string `json:"url"`
}
