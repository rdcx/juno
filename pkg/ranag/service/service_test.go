package service

import (
	"juno/pkg/api/client"
	"juno/pkg/api/node/dto"
	queryDto "juno/pkg/api/query/dto"
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

func TestQueryRange(t *testing.T) {
	t.Run("should combine results", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://node1.com:9090").
			Post("/query").
			Reply(200).
			JSON(queryDto.NewSuccessQueryResponse(
				map[string]string{
					"https://google.com/about": "Google About",
				},
			))

		gock.New("http://node2.com:9090").
			Post("/query").
			Reply(200).
			JSON(queryDto.NewSuccessQueryResponse(
				map[string]string{
					"https://microsoft.com":       "Microsoft",
					"https://microsoft.com/about": "Microsoft About",
				},
			))

		gock.New("http://node3.com:9090").
			Post("/query").
			Reply(200).
			JSON(queryDto.NewSuccessQueryResponse(
				map[string]string{
					"https://amazon.com/about": "Amazon About",
				},
			))

		svc := New(WithLogger(logrus.New()))

		svc.SetShards([shard.SHARDS][]string{
			0: {"node1.com:9090"},
			1: {"node2.com:9090"},
			2: {"node3.com:9090"},
		})

		values, err := svc.QueryRange(0, 3, &queryDto.Query{
			BasicQuery: &queryDto.BasicQuery{
				Title: &queryDto.StringMatch{
					Value:     "About",
					MatchType: queryDto.ContainsStringMatch,
				},
			},
		})

		if err == nil {
			t.Errorf("expected an error")
		}
	})
}
