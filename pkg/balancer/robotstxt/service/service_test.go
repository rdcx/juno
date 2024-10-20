package service

import (
	"juno/pkg/balancer/robotstxt/repo/mem"
	"testing"

	"github.com/h2non/gock"
)

func TestCanCrawlURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://example.com").
			Get("/robots.txt").
			Reply(200).
			BodyString("User-agent: *\nDisallow: /private")

		repo := mem.New()
		service := New(repo)

		canCrawl := service.CanCrawlURL("http://example.com")

		if !canCrawl {
			t.Errorf("expected to be able to crawl, got %v", canCrawl)
		}

		hit, err := repo.Get("example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if hit != "User-agent: *\nDisallow: /private" {
			t.Errorf("expected hit to be User-agent: *\nDisallow: /private, got %s", hit)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("no robots return true", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://example.com").
			Get("/robots.txt").
			Reply(404)

		repo := mem.New()
		service := New(repo)

		canCrawl := service.CanCrawlURL("http://example.com")

		if !canCrawl {
			t.Errorf("expected to be able to crawl, got %v", canCrawl)
		}

		hit, err := repo.Get("example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if hit != "" {
			t.Errorf("expected hit to be empty, got %s", hit)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})

	t.Run("robots return false", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://example.com").
			Get("/robots.txt").
			Reply(200).
			BodyString("User-agent: *\nDisallow: /")

		repo := mem.New()
		service := New(repo)

		canCrawl := service.CanCrawlURL("http://example.com")

		if canCrawl {
			t.Errorf("expected not to be able to crawl, got %v", canCrawl)
		}

		hit, err := repo.Get("example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if hit != "User-agent: *\nDisallow: /" {
			t.Errorf("expected hit to be User-agent: *\nDisallow: /, got %s", hit)
		}

		if !gock.IsDone() {
			t.Errorf("expected all mocks to be called")
		}
	})
}
