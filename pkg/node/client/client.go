package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"juno/pkg/node"
	domain "juno/pkg/node/crawl"
	crawlDto "juno/pkg/node/crawl/dto"
	extractionDto "juno/pkg/node/extraction/dto"
	"juno/pkg/util"
	"net/http"
)

func SendCrawlRequest(node string, url string) error {
	var req crawlDto.CrawlRequest

	req.URL = url

	b, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := http.Post("http://"+node+"/crawl", "application/json", bytes.NewBuffer(b))

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

func SendExtractionRequest(nodeAddr string, selectors []*extractionDto.Selector, fields []*extractionDto.Field) ([]map[string]interface{}, error) {
	b, err := json.Marshal(&extractionDto.ExtractionRequest{
		Selectors: selectors,
		Fields:    fields,
	})

	if err != nil {
		return nil, err
	}

	res, err := http.Post("http://"+nodeAddr+"/extraction", "application/json", bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, util.WrapErr(
			node.ErrFailedQueryRequest,
			fmt.Sprintf("status code: %d", res.StatusCode),
		)
	}

	var response extractionDto.ExtractionResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return response.Extractions, nil
}
