package service

import (
	"testing"

	balancerService "juno/pkg/node/balancer/service"
	fetcherService "juno/pkg/node/fetcher/service"
	pageRepo "juno/pkg/node/page/repo/mem"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"
)

func setupService(t *testing.T) *Service {
	balancerService := balancerService.New()
	pageRepo := pageRepo.New()
	pageService := pageService.New(pageRepo)
	storageService := storageService.New(
		t.TempDir(),
	)
	fetcherService := fetcherService.New()

	return New(
		balancerService,
		pageService,
		storageService,
		fetcherService,
	)
}

func TestCrawl(t *testing.T) {

}
