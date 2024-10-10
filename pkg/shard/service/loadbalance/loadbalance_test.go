package loadbalance

import (
	"fmt"
	"juno/pkg/shard/domain"
	"testing"

	"github.com/h2non/gock"
)

func makeTestNodes() [domain.SHARDS][]string {
	nodes := [domain.SHARDS][]string{}

	for i := 0; i < domain.SHARDS; i++ {
		nodes[i] = []string{fmt.Sprintf("http://node%d:8080", i)}
	}

	return nodes
}

func TestShardCrawl(t *testing.T) {
	t.Run("sends crawl request to random shard node", func(t *testing.T) {

		defer gock.Off()

		url := "http://example.com"

		// Set up mock responses for each shard node
		for i := 0; i < domain.SHARDS; i++ {
			nodeURL := fmt.Sprintf("http://node%d:8080", i)
			times := 0
			if domain.GetShard(url) == i {
				times = 1
			}
			gock.New(nodeURL).
				Post("/crawl").
				MatchType("json").
				JSON(map[string]string{"url": url}).
				Times(times)

		}

		// Create a new shard load balancer
		lb := NewLoadBalancer()
		lb.SetNodes(makeTestNodes())

		lb.Crawl(url)

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
