package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"juno/pkg/crawl/domain"
	"juno/pkg/util"
	"net/http"
)

func SendCrawlRequest(node string, url string) error {
	var req domain.CrawlRequest

	req.URL = url

	b, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := http.Post(node+"/crawl", "application/json", bytes.NewBuffer(b))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return util.WrapErr(
			domain.ErrFailedCrawlRequest,
			fmt.Sprintf("status code: %d", res.StatusCode),
		)
	}

	return nil
}
