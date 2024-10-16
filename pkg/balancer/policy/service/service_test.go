package service

import (
	"juno/pkg/balancer/policy"
	"juno/pkg/balancer/policy/repo/mem"
	"testing"
	"time"
)

func TestCanCrawl(t *testing.T) {
	t.Run("should return false when the last crawled is within the interval", func(t *testing.T) {
		svc := New(mem.New())

		p := &policy.CrawlPolicy{
			CrawlInterval: 10 * time.Minute,
			LastCrawled:   time.Now(),
		}

		canCrawl := svc.CanCrawl(p)

		if canCrawl {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("should return true when the last crawled is outside the interval", func(t *testing.T) {
		svc := New(mem.New())

		p := &policy.CrawlPolicy{
			CrawlInterval: 10 * time.Minute,
			LastCrawled:   time.Now().Add(-15 * time.Minute),
		}

		canCrawl := svc.CanCrawl(p)

		if !canCrawl {
			t.Errorf("expected true, got false")
		}
	})
}

func TestRecordCrawl(t *testing.T) {
	t.Run("should update the last crawled time and increment the times crawled", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		repo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   time.Now().Add(-10 * time.Minute),
			CrawlInterval: 5 * time.Second,
		})

		p, err := svc.Get("example.com")
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		err = svc.RecordCrawl(p)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		if p.LastCrawled.After(time.Now()) {
			t.Errorf("expected last crawled to be before now but got %v", p.LastCrawled)
		}

		if p.TimesCrawled != 1 {
			t.Errorf("expected times crawled to be 1 but got %d", p.TimesCrawled)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("should return policy when found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		repo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   time.Now().Add(-10 * time.Minute),
			CrawlInterval: 5 * time.Second,
		})

		p, err := svc.Get("example.com")

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if p == nil {
			t.Fatal("expected policy but got nil")
		}

		if p.Hostname != "example.com" {
			t.Errorf("expected example.com but got %s", p.Hostname)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("should set policy", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		repo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   time.Now().Add(-10 * time.Minute),
			CrawlInterval: 5 * time.Second,
		})

		p, err := svc.Get("example.com")

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if p == nil {
			t.Fatal("expected policy but got nil")
		}

		if p.Hostname != "example.com" {
			t.Errorf("expected example.com but got %s", p.Hostname)
		}
	})
}
