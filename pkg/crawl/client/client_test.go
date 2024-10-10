package client

import (
	"fmt"
	"strings"
	"testing"

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

		err := SendCrawlRequest("http://example.com", "http://shop.org")

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

				err := SendCrawlRequest("http://example.com", "http://shop.org")

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
