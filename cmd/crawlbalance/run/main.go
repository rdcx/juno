package main

import (
	"juno/pkg/balancer/crawl/handler"
	"juno/pkg/balancer/crawl/service"
	"juno/pkg/balancer/router"
)

func main() {
	balancerHandler := handler.New(
		service.New(),
	)
	router.New(balancerHandler)
}
