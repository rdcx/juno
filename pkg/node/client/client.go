package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"juno/pkg/node"
	domain "juno/pkg/node/crawl"
	crawlDto "juno/pkg/node/crawl/dto"
	extractionDto "juno/pkg/node/extraction/dto"
	infoDto "juno/pkg/node/info/dto"
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

func SendExtractionRequest(nodeAddr string, shard int, selectors []*extractionDto.Selector, fields []*extractionDto.Field) ([]map[string]interface{}, error) {
	b, err := json.Marshal(&extractionDto.ExtractionRequest{
		Shard:     shard,
		Selectors: selectors,
		Fields:    fields,
	})

	if err != nil {
		return nil, err
	}

	res, err := http.Post("http://"+nodeAddr+"/extract", "application/json", bytes.NewBuffer(b))

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

func SendInfoRequest(nodeAddr string) (*infoDto.InfoResponse, error) {
	res, err := http.Get("http://" + nodeAddr + "/info")

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, util.WrapErr(
			node.ErrFailedInfoRequest,
			fmt.Sprintf("status code: %d", res.StatusCode),
		)
	}

	var response infoDto.InfoResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
