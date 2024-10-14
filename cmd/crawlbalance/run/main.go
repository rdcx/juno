package main

import (
	"juno/pkg/crawlbalance/balanacer/service"
	"juno/pkg/crawlbalance/balancer/handler"
	"juno/pkg/crawlbalance/router"
)

func main() {
	balancerHandler := handler.New(
		service.NewLoadBalancer(),
	)
	router.New(balancerHandler)
}
