package client

import (
	"fmt"
	"strings"
	"testing"

	extractionDto "juno/pkg/node/extraction/dto"
	"juno/pkg/node/info"
	infoDto "juno/pkg/node/info/dto"

	"github.com/h2non/gock"
)

func TestSendCrawlRequest(t *testing.T) {

	t.Run("sends crawl request", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://example.com").
			Post("/crawl").
			MatchType("json").
			JSON(map[string]string{"url": "http://shop.org"}).
			Times(1).
			Reply(200)

		err := SendCrawlRequest("example.com", "http://shop.org")

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})

	t.Run("returns error if request fails", func(t *testing.T) {

		tests := []struct {
			statusCode int
		}{
			{500},
			{404},
			{429},
			{400},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("status code %d", test.statusCode), func(t *testing.T) {

				defer gock.Off()

				gock.New("http://example.com").
					Post("/crawl").
					MatchType("json").
					JSON(map[string]string{"url": "http://shop.org"}).
					Times(1).
					Reply(test.statusCode)

				err := SendCrawlRequest("example.com", "http://shop.org")

				if err == nil {
					t.Errorf("Expected an error")
				}

				if !strings.Contains(err.Error(), fmt.Sprintf("status code: %d", test.statusCode)) {
					t.Errorf("Unexpected error: %s", err)
				}

				if !gock.IsDone() {
					t.Errorf("Not all expectations were met")
				}
			})
		}
	})
}

func TestSendExtractionRequest(t *testing.T) {

	t.Run("sends extraction request", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://node1.com:8080").
			Post("/extract").
			JSON(map[string]interface{}{
				"shard": 0,
				"selectors": []map[string]string{
					{"id": "1", "value": "#productTitle"},
				},
				"fields": []map[string]string{
					{"id": "1", "selector_id": "1", "name": "product_title"},
				},
			}).
			Times(1).
			Reply(200).
			JSON(extractionDto.NewSuccessExtractionResponse(
				[]map[string]interface{}{
					{"product_title": "charger"},
				},
			))

		res, err := SendExtractionRequest("node1.com:8080",
			0,
			[]*extractionDto.Selector{
				{
					ID:    "1",
					Value: "#productTitle",
				},
			},
			[]*extractionDto.Field{
				{
					ID:         "1",
					SelectorID: "1",
					Name:       "product_title",
				},
			},
		)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if res == nil {
			t.Errorf("Expected non-nil response")
		}

		if res[0]["product_title"] != "charger" {
			t.Errorf("Expected charger, got %s", res[0]["product_title"])
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}

func TestSendInfoRequest(t *testing.T) {

	t.Run("sends info request", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://node1.com:8080").
			Get("/info").
			Times(1).
			Reply(200).
			JSON(infoDto.NewSuccessInfoResponse(
				&info.Info{
					PageCount: 100,
				}),
			)

		res, err := SendInfoRequest("node1.com:8080")

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if res == nil {
			t.Fatal("Expected non-nil response")
		}

		if res.Info == nil {
			t.Fatal("Expected non-nil info")
		}

		if res.Info.PageCount != 100 {
			t.Errorf("Expected 100, got %d", res.Info.PageCount)
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
