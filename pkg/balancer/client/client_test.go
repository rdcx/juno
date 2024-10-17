package client

import (
	"testing"

	"github.com/h2non/gock"
)

func TestCrawlURLs(t *testing.T) {
	t.Run("should make crawl urls request", func(t *testing.T) {
		defer gock.Off()

		baseURL := "http://localhost:8080"

		gock.New(baseURL).
			Post("/crawl/urls").
			JSON(map[string][]string{"urls": {"http://example.com"}}).
			Reply(200)

		client := New(baseURL)

		err := client.CrawlURLs([]string{"http://example.com"})

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
			Post("/crawl/urls").
			JSON(map[string][]string{"urls": {"http://example.com"}}).
			Reply(500)

		client := New(baseURL)

		err := client.CrawlURLs([]string{"http://example.com"})

		if err == nil {
			t.Errorf("Expected an error")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}

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
