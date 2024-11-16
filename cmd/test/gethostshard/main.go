package main

import (
	"fmt"
	"juno/pkg/shard"
)

func main() {
	shard := shard.GetShard("https://en.wikipedia.org/wiki/2024_Japanese_general_election")

	fmt.Printf("Shard: %d\n", shard)
}
