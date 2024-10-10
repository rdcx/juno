package client

import (
	"testing"

	"github.com/h2non/gock"
)

func TestSendCrawlRequest(t *testing.T) {
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
}
