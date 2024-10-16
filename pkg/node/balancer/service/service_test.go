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

func TestWithBalancerFetchInterval(t *testing.T) {
	t.Run("should panic when api client is not set", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		WithBalancerFetchInterval(0)(&Service{})
	})

	t.Run("should fetch balancers every interval", func(t *testing.T) {
		defer gock.Off()

		baseURL := "http://localhost:8080"
		client := client.New(baseURL)

		gock.New(baseURL).
			Get("/balancers").
			Reply(200).
			JSON(dto.AllShardsNodesResponse{
				Shards: map[int][]string{
					1: {"balancer1.com:9090", "balancer2.com:9090"},
				},
			})

		svc := New(WithLogger(logrus.New()),
			WithApiClient(client))
		WithBalancerFetchInterval(50 * time.Millisecond)(svc)

		time.Sleep(75 * time.Millisecond)

		if len(svc.balancers) == 0 {
			t.Errorf("expected shards to be fetched")
		}

		gock.New(baseURL).
			Get("/balancers").
			Reply(200).
			JSON(dto.AllShardsNodesResponse{
				Shards: map[int][]string{
					1: {"balancer3.com:8080", "balancer4.com:8080"},
				},
			})

		time.Sleep(75 * time.Millisecond)

		if len(svc.balancers) == 0 {
			t.Errorf("expected shards to be fetched")
		}

		if svc.balancers[1][0] != "balancer3.com:8080" {
			t.Errorf("expected balancer3 but got %s", svc.balancers[1][0])
		}

		if svc.balancers[1][1] != "balancer4.com:8080" {
			t.Errorf("expected balanacer4 but got %s", svc.balancers[1][1])
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}

func TestSendCrawlRequest(t *testing.T) {
	t.Run("should send crawl request", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://balancer1.com:9090").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Reply(200)

		svc := New(WithLogger(logrus.New()))

		svc.SetBalancers([shard.SHARDS][]string{
			72435: {"balancer1.com:9090"},
		})

		err := svc.SendCrawlRequest("http://example.com")

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})

	t.Run("should return error on failed request", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://balancer1.com:9090").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com"}).
			Reply(500)

		svc := New(WithLogger(logrus.New()))

		svc.SetBalancers([shard.SHARDS][]string{
			72435: {"balancer1.com:9090"},
		})

		err := svc.SendCrawlRequest("http://example.com")

		if err == nil {
			t.Errorf("Expected an error")
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
