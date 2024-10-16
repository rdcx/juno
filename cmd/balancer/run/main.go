package main

import (
	"flag"
	"juno/pkg/api/client"
	"juno/pkg/balancer/crawl/handler"
	"juno/pkg/balancer/crawl/service"
	"juno/pkg/balancer/router"
	"time"
)

func main() {

	var apiFlag string
	flag.StringVar(&apiFlag, "api", "http://localhost:8080", "API URL")

	flag.Parse()

	if apiFlag == "" {
		panic("API URL is required")
	}

	apiClient := client.New(apiFlag)

	balancerHandler := handler.New(
		service.New(
			service.WithApiClient(apiClient),
			service.WithShardFetchInterval(time.Minute*5),
		),
	)
	router.New(balancerHandler)
}
