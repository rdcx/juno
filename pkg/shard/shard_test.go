package shard

import (
	"fmt"
	"testing"
)

func TestGetShard(t *testing.T) {
	t.Run("returns a shard number", func(t *testing.T) {
		host := "shop.com"
		shard := GetShard(host)

		if shard < 0 || shard >= SHARDS {
			t.Errorf("Expected shard to be between 0 and %d, got %d", SHARDS, shard)
		}
	})

	t.Run("distribution is uniform", func(t *testing.T) {
		try := 100_000_000
		counts := make(map[int]int)

		for i := 0; i < try; i++ {
			host := fmt.Sprintf("%dsh%dop%d.com", i, i, i)
			shard := GetShard(host)
			counts[shard]++
		}

		// check each shard has at least 1% of the total
		for shard, count := range counts {
			if count < try/SHARDS/100 {
				t.Errorf("Expected shard %d to have at least %d, got %d", shard, try/SHARDS/100, count)
			}
		}
	})
}
