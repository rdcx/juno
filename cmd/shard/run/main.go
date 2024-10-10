package main

import (
	"juno/pkg/shard/handler"
	"juno/pkg/shard/router"
	"juno/pkg/shard/service/loadbalance"
)

func main() {
	shardLb := handler.NewLoadBalanceHandler(
		loadbalance.NewLoadBalancer(),
	)
	router.RunShardService(shardLb, "8080")
}
