package client

import (
	"bytes"
	"encoding/json"
	"juno/pkg/crawl/domain"
	"net/http"
)

func SendCrawlRequest(node string, url string) error {
	var req domain.CrawlRequest

	req.URL = url

	b, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = http.Post(node+"/crawl", "application/json", bytes.NewBuffer(b))

	if err != nil {
		return err
	}

	return nil
}
