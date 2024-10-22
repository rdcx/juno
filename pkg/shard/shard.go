package shard

import (
	"errors"

	"github.com/spaolacci/murmur3"
)

var ErrInvalidTotal = errors.New("total must be greater than 0")

const SHARDS = 100_000

func GetShard(hostname string) int {
	hash := murmur3.Sum32([]byte(hostname))
	return int(hash % uint32(SHARDS))
}

func GetShardRange(offset, total int) ([]int, error) {

	if offset < 0 {
		return nil, ErrInvalidTotal
	}
	if total < 1 {
		return nil, ErrInvalidTotal
	}

	shards := []int{}

	for i := offset; i < offset+total; i++ {
		shards = append(shards, i%SHARDS)
	}

	return shards, nil
}
