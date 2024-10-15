package shard

import (
	"github.com/spaolacci/murmur3"
)

const SHARDS = 100_000
const BLOCK_SIZE = 1_000
const BLOCKS = SHARDS / BLOCK_SIZE

func GetShard(hostname string) int {
	hash := murmur3.Sum32([]byte(hostname))
	return int(hash % uint32(SHARDS))
}

var example = "ABA2-34"

var another = "[0,1000], [1000,2000], [2000,3000]"

type Range struct {
	ID   string
	From int
	To   int
}

var ranges = []Range{}

const FORMAT = "A-Z0-9A-Z0-9A-Z"

type ShardAssignment string

func NewShardAssignment(addr string) ShardAssignment {
	if len(addr) != 5 {
		panic("invalid address")
	}

	for i, c := range addr {
		if i%2 == 0 {
			if c < 'A' || c > 'Z' {
				panic("invalid address")
			}
		} else {
			if c < '0' || c > '9' {
				panic("invalid address")
			}
		}
	}
	return ShardAssignment(addr)
}
