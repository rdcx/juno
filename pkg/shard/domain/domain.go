package domain

import "github.com/spaolacci/murmur3"

const SHARDS = 100

func GetShard(hostname string) int {
	hash := murmur3.Sum32([]byte(hostname))
	return int(hash % uint32(SHARDS))
}
