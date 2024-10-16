package client

import (
	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("should make crawl request", func(t *testing.T) {
		defer gock.Off()

		baseURL := "http://localhost:8080"

		gock.New(baseURL).
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Reply(200)

		client := New(baseURL)

		err := client.Crawl("http://example.com")

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})

	t.Run("should return error on failed request", func(t *testing.T) {
		defer gock.Off()

		baseURL := "http://localhost:8080"

		gock.New(baseURL).
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Reply(500)

		client := New(baseURL)

		err := client.Crawl("http://example.com")

		if err == nil {
			t.Errorf("Expected an error")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
