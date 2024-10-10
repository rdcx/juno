package service

import (
	"testing"
	"time"
)

func TestPolicyManager(t *testing.T) {
	t.Run("can crawl when last crawl is before interval", func(t *testing.T) {
		pm := NewPolicyManager()
		pm.SetPolicy("example.com", &Policy{
			Interval:  time.Hour,
			LastCrawl: time.Now().Add(-time.Hour),
		})

		if !pm.CanCrawl("example.com") {
			t.Error("expected to be able to crawl")
		}
	})
}
