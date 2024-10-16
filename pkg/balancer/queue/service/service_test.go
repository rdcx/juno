package service

import (
	"context"
	"juno/pkg/balancer/policy"
	polRepo "juno/pkg/balancer/policy/repo/mem"
	"juno/pkg/balancer/policy/service"
	"juno/pkg/balancer/queue"
	"juno/pkg/balancer/queue/repo/mem"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type testCrawlService struct {
	returnErr error
	hits      []string
}

func (s *testCrawlService) Crawl(url string) error {
	if s.returnErr != nil {
		return s.returnErr
	}

	s.hits = append(s.hits, url)
	return nil
}

func TestProcess(t *testing.T) {
	t.Run("should process queue", func(t *testing.T) {
		// Given
		crawlSvc := &testCrawlService{}
		repo := mem.New()
		polSvc := service.New(polRepo.New())
		svc := New(logrus.New(), repo, crawlSvc, polSvc)
		repo.Push("http://example.com")

		ctx, cancel := context.WithCancel(context.Background())

		// Run Process in a separate goroutine
		go func() {
			err := svc.Process(ctx)
			if err != nil && err != queue.ErrProcessQueueCancelled {
				t.Errorf("expected no error but got %v", err)
			}
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		time.Sleep(100 * time.Millisecond) // Give time for cancellation to propagate

		// Then
		if len(crawlSvc.hits) != 1 {
			t.Errorf("expected 1 hit but got %d", len(crawlSvc.hits))
		}

		if crawlSvc.hits[0] != "http://example.com" {
			t.Errorf("expected http://example.com but got %s", crawlSvc.hits[0])
		}

		// check policy was set
		p, err := polSvc.Get("example.com")

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if p == nil {
			t.Fatal("expected policy but got nil")
		}

		if p.Hostname != "example.com" {
			t.Errorf("expected example.com but got %s", p.Hostname)
		}

		if p.TimesCrawled != 1 {
			t.Errorf("expected 1 but got %d", p.TimesCrawled)
		}

		if p.LastCrawled.IsZero() {
			t.Error("expected LastCrawled to be set")
		}
	})

	t.Run("does not process url when policy says it can't", func(t *testing.T) {
		// Given
		crawlSvc := &testCrawlService{}
		repo := mem.New()
		polRepo := polRepo.New()
		polSvc := service.New(polRepo)
		svc := New(logrus.New(), repo, crawlSvc, polSvc)
		repo.Push("http://example.com")

		lastCrawledTime := time.Now()
		polRepo.Set("example.com", &policy.CrawlPolicy{
			Hostname:      "example.com",
			LastCrawled:   lastCrawledTime,
			CrawlInterval: 1 * time.Hour,
		})

		ctx, cancel := context.WithCancel(context.Background())

		// Run Process in a separate goroutine
		go func() {
			err := svc.Process(ctx)
			if err != nil && err != queue.ErrProcessQueueCancelled {
				t.Errorf("expected no error but got %v", err)
			}
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		time.Sleep(100 * time.Millisecond) // Give time for cancellation to propagate

		// Then
		if len(crawlSvc.hits) != 0 {
			t.Errorf("expected 0 hits but got %d", len(crawlSvc.hits))
		}

		// check policy was set
		p, err := polSvc.Get("example.com")

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if p == nil {
			t.Fatal("expected policy but got nil")
		}

		if p.Hostname != "example.com" {
			t.Errorf("expected example.com but got %s", p.Hostname)
		}

		if p.TimesCrawled != 0 {
			t.Errorf("expected 0 but got %d", p.TimesCrawled)
		}

		if !p.LastCrawled.Equal(lastCrawledTime) {
			t.Error("expected LastCrawled to be unchanged")
		}
	})
}
