package main

import (
	"flag"
	"juno/pkg/api/client"
	balancerService "juno/pkg/node/balancer/service"
	crawlHandler "juno/pkg/node/crawl/handler"
	crawlService "juno/pkg/node/crawl/service"
	fetcherService "juno/pkg/node/fetcher/service"
	htmlService "juno/pkg/node/html/service"
	pageRepo "juno/pkg/node/page/repo/bolt"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"
	"time"

	"juno/pkg/node/router"

	"github.com/sirupsen/logrus"
)

func main() {

	var apiURL string
	flag.StringVar(&apiURL, "api-url", "http://localhost:8080", "URL of the API server")
	var pageDBPath string
	flag.StringVar(&pageDBPath, "page-db-path", "page.db", "Path to the page database")
	var storageDir string
	flag.StringVar(&storageDir, "storage-dir", "storage", "Directory to store downloaded HTML files")

	flag.Parse()

	if apiURL == "" {
		panic("api-url flag is required")
	}

	logger := logrus.New()

	pageRepo, err := pageRepo.New(pageDBPath)

	if err != nil {
		panic(err)
	}

	storageService := storageService.New(storageDir)
	pageService := pageService.New(pageRepo)

	htmlService := htmlService.New()

	balancerService := balancerService.New(
		balancerService.WithApiClient(
			client.New(apiURL),
		),
		balancerService.WithBalancerFetchInterval(time.Minute),
		balancerService.WithLogger(logger),
	)

	fetcherService := fetcherService.New()

	crawlService := crawlService.New(
		balancerService,
		pageService,
		storageService,
		fetcherService,
		htmlService,
	)

	crawlHandler := crawlHandler.New(logger, crawlService)
	r := router.New(
		crawlHandler,
	)

	r.Run(":8081")
}
