package service

import (
	"context"
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
		svc := New(logrus.New(), repo, crawlSvc)
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
	})
}
