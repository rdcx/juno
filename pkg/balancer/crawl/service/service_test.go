package service

import (
	"juno/pkg/api/client"
	"juno/pkg/api/node/dto"
	"juno/pkg/shard"
	"testing"
	"time"

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
