package loadbalance

import (
	"juno/pkg/shard"
	"testing"

	"github.com/h2non/gock"
)

func TestCrawl(t *testing.T) {
	t.Run("sends crawl request to random shard node", func(t *testing.T) {

		defer gock.Off()

		url := "http://example.com"

		gock.New("http://node1.com:8080").
			Post("/crawl").
			MatchType("json").
			JSON(map[string]string{"url": url})

		svc := New()
		svc.SetNodes([shard.SHARDS][]string{
			72435: {"http://node1.com:8080"},
		})

		svc.Crawl(url)

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
