package main

import (
	"encoding/json"
	"juno/pkg/api/node/dto"
	"net/http"
	"os"
	"strconv"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		panic("Please provide a shard ID")
	}

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

	// convert string to int
	conv, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	nodes, ok := res.Shards[conv]
	if !ok {
		println("No nodes found for shard", conv)
		return
	}

	for _, node := range nodes {
		println(node)
	}
}
