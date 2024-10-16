package service

import (
	"context"
	"juno/pkg/api/client"
	"juno/pkg/api/node/dto"
	"juno/pkg/balancer/policy"
	"juno/pkg/balancer/queue"
	"juno/pkg/balancer/queue/repo/mem"
	"juno/pkg/shard"
	"testing"
	"time"

	polRepo "juno/pkg/balancer/policy/repo/mem"
	polService "juno/pkg/balancer/policy/service"

	queueRepo "juno/pkg/balancer/queue/repo/mem"
	queueService "juno/pkg/balancer/queue/service"

	"github.com/h2non/gock"
	"github.com/sirupsen/logrus"
)

func TestWithShardFetchInterval(t *testing.T) {
	t.Run("should panic when api client is not set", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		WithShardFetchInterval(0)(&Service{})
	})

	t.Run("should fetch shards every interval", func(t *testing.T) {
		defer gock.Off()

		baseURL := "http://localhost:8080"
		client := client.New(baseURL)

		gock.New(baseURL).
			Get("/shards").
			Reply(200).
			JSON(dto.AllShardsNodesResponse{
				Shards: map[int][]string{
					1: {"node1.com:9090", "node2.com:9090"},
				},
			})

		svc := New(WithLogger(logrus.New()),
			WithApiClient(client))
		WithShardFetchInterval(50 * time.Millisecond)(svc)

		time.Sleep(75 * time.Millisecond)

		if len(svc.shards) == 0 {
			t.Errorf("expected shards to be fetched")
		}

		gock.New(baseURL).
			Get("/shards").
			Reply(200).
			JSON(dto.AllShardsNodesResponse{
				Shards: map[int][]string{
					1: {"node3.com:8080", "node4.com:8080"},
				},
			})

		time.Sleep(75 * time.Millisecond)

		if len(svc.shards) == 0 {
			t.Errorf("expected shards to be fetched")
		}

		if svc.shards[1][0] != "node3.com:8080" {
			t.Errorf("expected node3 but got %s", svc.shards[1][0])
		}

		if svc.shards[1][1] != "node4.com:8080" {
			t.Errorf("expected node4 but got %s", svc.shards[1][1])
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}

func TestCrawl(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Times(1).
			Reply(200).
			JSON(map[string]string{"message": "ok"})

		svc := New(WithLogger(logrus.New()))
		svc.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})

		svc.Crawl("http://example.com")

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})

	t.Run("tries 3 times before giving up", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Times(3).
			Reply(500).
			JSON(map[string]string{"message": "internal server error"})

		svc := New(WithLogger(logrus.New()))
		svc.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})

		err := svc.Crawl("http://example.com")

		if err == nil {
			t.Errorf("expected an error")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}

func TestProcess(t *testing.T) {
	t.Run("should process queue", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Times(1).
			Reply(200).
			JSON(map[string]string{"message": "ok"})

		logger := logrus.New()
		// Given
		queueRepo := mem.New()
		polSvc := polService.New(polRepo.New())
		queueService := queueService.New(logger, queueRepo)
		crawlService := New(
			WithLogger(logger),
			WithQueueService(queueService),
			WithPolicyService(polSvc),
		)
		crawlService.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})

		pol := policy.New("example.com")
		lastCrawledTime := time.Now().Add(-30 * time.Minute)
		pol.LastCrawled = lastCrawledTime
		polSvc.Set("example.com", pol)

		queueRepo.Push("http://example.com")

		ctx, cancel := context.WithCancel(context.Background())

		// Run Process in a separate goroutine
		go func() {
			err := crawlService.ProcessQueue(ctx)
			if err != nil && err != queue.ErrProcessQueueCancelled {
				t.Errorf("expected no error but got %v", err)
			}
		}()

		time.Sleep(49 * time.Millisecond)
		cancel()

		time.Sleep(99 * time.Millisecond) // Give time for cancellation to propagate

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

		if p.LastCrawled.Equal(lastCrawledTime) {
			t.Error("expected LastCrawled to be updated")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})

	t.Run("does not process url when policy says it can't", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/crawl").
			Times(0)

		logger := logrus.New()
		queueRepo := queueRepo.New()
		polSvc := polService.New(polRepo.New())
		queueService := queueService.New(logger, queueRepo)
		crawlService := New(
			WithLogger(logger),
			WithQueueService(queueService),
			WithPolicyService(polSvc),
		)
		crawlService.SetShards([shard.SHARDS][]string{
			72435: {"node1.com:9090"},
		})

		lastCrawledTime := time.Now()
		pol := policy.New("example.com")
		// simulate host has just been crawled
		pol.LastCrawled = lastCrawledTime
		pol.TimesCrawled = 1
		polSvc.Set("example.com", pol)

		queueRepo.Push("http://example.com")

		ctx, cancel := context.WithCancel(context.Background())

		// Run Process in a separate goroutine
		go func() {
			err := crawlService.ProcessQueue(ctx)
			if err != nil && err != queue.ErrProcessQueueCancelled {
				t.Errorf("expected no error but got %v", err)
			}
		}()

		time.Sleep(49 * time.Millisecond)
		cancel()

		time.Sleep(99 * time.Millisecond) // Give time for cancellation to propagate

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

		if !p.LastCrawled.Equal(lastCrawledTime) {
			t.Error("expected LastCrawled to be unchanged")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
