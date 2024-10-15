package main

import (
	"encoding/json"
	"juno/pkg/api/node/dto"
	"net/http"
)

func main() {

	resp, err := http.Get("http://localhost:8080/shards")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var res dto.AllShardsNodesResponse

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		panic(err)
	}

	shards := 100_000

	for i := 0; i < shards; i++ {
		_, ok := res.Shards[i]
		if !ok {
			println("No nodes found for shard", i)
		}
	}
}
