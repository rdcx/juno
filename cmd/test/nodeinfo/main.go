package main

import (
	"fmt"
	"juno/pkg/api/client"

	nclient "juno/pkg/node/client"
)

func uniqueAddresses(list []string) []string {
	unique := make(map[string]bool)
	var addresses []string

	for _, addr := range list {
		if !unique[addr] {
			unique[addr] = true
			addresses = append(addresses, addr)
		}
	}

	return addresses
}

func main() {

	apiClient := client.New("http://localhost:8080")

	shards, err := apiClient.GetShards()

	if err != nil {
		panic(err)
	}

	totalPages := 0

	var nodeAddrs []string

	for _, shard := range shards.Shards {
		nodeAddrs = append(nodeAddrs, shard...)
	}

	nodeAddrs = uniqueAddresses(nodeAddrs)

	for _, addr := range nodeAddrs {
		res, err := nclient.SendInfoRequest(addr)

		if err != nil {
			fmt.Printf("failed to get info from %s: %v\n", addr, err)
			continue
		}

		totalPages += res.Info.PageCount
	}

	fmt.Println("Total pages:", totalPages)
}
