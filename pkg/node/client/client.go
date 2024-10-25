package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	domain "juno/pkg/node/crawl"
	"juno/pkg/node/crawl/dto"
	"juno/pkg/util"
	"net/http"
)

func SendCrawlRequest(node string, url string) error {
	var req dto.CrawlRequest

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

// func SendExtractionRequest(nodeAddr string, req extractorDto.ExtractionRequest) (interface{}, error) {
// 	b, err := json.Marshal(req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := http.Post("http://"+nodeAddr+"/extraction", "application/json", bytes.NewBuffer(b))

// 	if err != nil {
// 		return nil, err
// 	}

// 	if res.StatusCode != http.StatusOK {
// 		return nil, util.WrapErr(
// 			node.ErrFailedQueryRequest,
// 			fmt.Sprintf("status code: %d", res.StatusCode),
// 		)
// 	}

// 	var response interface{}

// 	err = json.NewDecoder(res.Body).Decode(&response)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }
