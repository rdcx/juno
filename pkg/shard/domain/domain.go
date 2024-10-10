package domain

import (
	"errors"

	"github.com/spaolacci/murmur3"
)

var ErrNoNodesAvailableInShard = errors.New("no nodes available in shard")

const SHARDS = 100

func GetShard(hostname string) int {
	hash := murmur3.Sum32([]byte(hostname))
	return int(hash % uint32(SHARDS))
}
