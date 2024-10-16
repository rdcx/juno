package service

import (
	"context"
	"errors"
	"juno/pkg/node/crawl"
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
)

const clothesPage = `<html>
	<head>
		<title>Shop</title>
	</head>
	<body>
		<h1>Clothes</h1>
	</body>
</html>`

func TestFetchPage(t *testing.T) {
	t.Run("fetches a page", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			Reply(200).
			BodyString(clothesPage)

		body, status, finalURL, err := FetchPage(context.Background(), "https://shop.com/clothes")

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if status != 200 {
			t.Errorf("Expected status to be 200, got %d", status)
		}

		if finalURL != "https://shop.com/clothes" {
			t.Errorf("Expected finalURL to be https://shop.com/clothes, got %s", finalURL)
		}

		if string(body) != clothesPage {
			t.Errorf("Expected page to be %s, got %s", clothesPage, body)
		}
	})

	t.Run("handles 4XX client error", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			Reply(400)

		body, status, finalURL, err := FetchPage(context.Background(), "https://shop.com/clothes")

		if !errors.Is(err, crawl.Err400) {
			t.Errorf("Expected error to be Err400, got %s", err)
		}

		if status != 400 {
			t.Errorf("Expected status to be 400, got %d", status)
		}

		if finalURL != "https://shop.com/clothes" {
			t.Errorf("Expected finalURL to be https://shop.com/clothes, got %s", finalURL)
		}

		if body != nil {
			t.Errorf("Expected page to be nil, got %s", body)
		}
	})

	t.Run("handles 404 page not found", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			Reply(404)

		body, status, finalURL, err := FetchPage(context.Background(), "https://shop.com/clothes")

		if err == nil {
			t.Error("Expected an error, got nil")
		}

		if status != 404 {
			t.Errorf("Expected status to be 404, got %d", status)
		}

		if finalURL != "https://shop.com/clothes" {
			t.Errorf("Expected finalURL to be https://shop.com/clothes, got %s", finalURL)
		}

		if body != nil {
			t.Errorf("Expected page to be nil, got %s", body)
		}
	})

	t.Run("handles network errors", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			ReplyError(errors.New("network error"))

		body, status, finalURL, err := FetchPage(context.Background(), "https://shop.com/clothes")

		if !strings.Contains(err.Error(), "network error") {
			t.Errorf("Expected error to contain 'network error', got %s", err)
		}

		if status != 0 {
			t.Errorf("Expected status to be 0, got %d", status)
		}

		if finalURL != "" {
			t.Errorf("Expected finalURL to be empty, got %s", finalURL)
		}

		if body != nil {
			t.Errorf("Expected page to be nil, got %s", body)
		}
	})

	t.Run("handles 429 rate limit error", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			Reply(429)

		body, status, finalURL, err := FetchPage(context.Background(), "https://shop.com/clothes")

		if err != crawl.Err429 {
			t.Errorf("Expected error to be Err429, got %s", err)
		}

		if status != 429 {
			t.Errorf("Expected status to be 429, got %d", status)
		}

		if finalURL != "https://shop.com/clothes" {
			t.Errorf("Expected finalURL to be https://shop.com/clothes, got %s", finalURL)
		}

		if body != nil {
			t.Errorf("Expected page to be nil, got %s", body)
		}
	})

	t.Run("handles timeout", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://shop.com").
			Get("/clothes").
			Response.ResponseDelay = time.Millisecond * 75

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()

		body, status, finalURL, err := FetchPage(ctx, "https://shop.com/clothes")

		if err != crawl.ErrContextDone {
			t.Errorf("Expected error to be ErrContextDone, got %s", err)
		}

		if status != 0 {
			t.Errorf("Expected status to be 0, got %d", status)
		}

		if finalURL != "" {
			t.Errorf("Expected finalURL to be empty, got %s", finalURL)
		}

		if body != nil {
			t.Errorf("Expected page to be nil, got %s", body)
		}
	})
}
