package shard

import "testing"

func TestGetShard(t *testing.T) {
	example := "example.com"

	t.Run("should return shard", func(t *testing.T) {
		shard := GetShard(example)

		if shard != 72435 {
			t.Errorf("expected shard 72435, got %d", shard)
		}
	})
}

func TestGetShardRange(t *testing.T) {
	t.Run("should return error when offset is less than 0", func(t *testing.T) {
		_, err := GetShardRange(-1, 1)

		if err != ErrInvalidTotal {
			t.Errorf("expected error to be returned")
		}
	})

	t.Run("should return error when total is less than 1", func(t *testing.T) {
		_, err := GetShardRange(0, 0)

		if err != ErrInvalidTotal {
			t.Errorf("expected error to be returned")
		}
	})

	t.Run("should return shard range", func(t *testing.T) {
		shards, err := GetShardRange(0, 10)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(shards) != 10 {
			t.Errorf("expected 10 shards, got %d", len(shards))
		}

		for i, shard := range shards {
			if shard != i {
				t.Errorf("expected shard %d, got %d", i, shard)
			}
		}
	})
}
