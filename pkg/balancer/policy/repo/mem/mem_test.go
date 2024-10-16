package mem

import (
	"juno/pkg/balancer/policy"
	"testing"
	"time"
)

func TestRepo(t *testing.T) {
	t.Run("should return nil when policy is not found", func(t *testing.T) {
		repo := New()

		p, err := repo.Get("example.com")

		if err != policy.ErrPolicyNotFound {
			t.Errorf("expected no error but got %v", err)
		}

		if p != nil {
			t.Errorf("expected nil but got %v", p)
		}
	})

	t.Run("should return policy when found", func(t *testing.T) {
		repo := New()

		repo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   time.Now().Add(-10 * time.Minute),
			CrawlInterval: 5 * time.Second,
		})

		p, err := repo.Get("example.com")

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
		repo := New()

		repo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   time.Now().Add(-10 * time.Minute),
			CrawlInterval: 5 * time.Second,
		})

		p, err := repo.Get("example.com")

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
