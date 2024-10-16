package main

import (
	"flag"
	"juno/pkg/api/client"

	queueRepo "juno/pkg/balancer/queue/repo/bolt"
	queueService "juno/pkg/balancer/queue/service"

	policyRepo "juno/pkg/balancer/policy/repo/bolt"
	policyService "juno/pkg/balancer/policy/service"

	crawlHandler "juno/pkg/balancer/crawl/handler"
	crawlService "juno/pkg/balancer/crawl/service"

	"juno/pkg/balancer/router"

	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()

	var apiFlag string
	flag.StringVar(&apiFlag, "api", "http://localhost:8080", "API URL")

	var policyDBPath string
	flag.StringVar(&policyDBPath, "policy-db", "policy.db", "Policy DB Path")

	var queueDBPath string
	flag.StringVar(&queueDBPath, "queue-db", "queue.db", "Queue DB Path")

	flag.Parse()

	if apiFlag == "" {
		panic("API URL is required")
	}

	apiClient := client.New(apiFlag)

	policyRepo, err := policyRepo.New(policyDBPath)

	if err != nil {
		panic(err)
	}

	policyService := policyService.New(policyRepo)

	queueRepo, err := queueRepo.New(queueDBPath)

	if err != nil {
		panic(err)
	}

	queueService := queueService.New(
		logger,
		queueRepo,
	)

	crawlService := crawlService.New(
		crawlService.WithLogger(logger),
		crawlService.WithApiClient(apiClient),
		crawlService.WithQueueService(queueService),
		crawlService.WithPolicyService(policyService),
	)

	crawlHandler := crawlHandler.New(
		logger,
		queueService,
	)

	go func() {

	}()

	r := router.New(crawlHandler)

	r.Run(":8080")
}
