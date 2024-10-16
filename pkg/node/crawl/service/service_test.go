package service

import (
	"context"
	"testing"

	balancerService "juno/pkg/node/balancer/service"
	fetcherService "juno/pkg/node/fetcher/service"
	htmlService "juno/pkg/node/html/service"
	"juno/pkg/node/page"
	pageRepo "juno/pkg/node/page/repo/mem"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"
	"juno/pkg/shard"

	"github.com/h2non/gock"
)

func setupService(t *testing.T) *Service {
	htmlService := htmlService.New()
	balancerService := balancerService.New()
	pageRepo := pageRepo.New()
	pageService := pageService.New(pageRepo)
	storageService := storageService.New(
		t.TempDir(),
	)
	fetcherService := fetcherService.New()

	balancerService.SetBalancers([shard.SHARDS][]string{
		72435: {"balancer1:8080"},
	})

	return New(
		balancerService,
		pageService,
		storageService,
		fetcherService,
		htmlService,
	)
}

var testFile = []byte(`<html>
	<head>
		<title>Test</title>
	</head>
	<body>
		<a href="http://example.com/about">Example</a>
	</body>
</html>`)

func TestCrawl(t *testing.T) {
	t.Run("should crawl a page", func(t *testing.T) {
		s := setupService(t)

		defer gock.Off()

		gock.New("http://example.com").
			Get("/home").
			Reply(200).
			BodyString(string(testFile))

		gock.New("http://balancer1:8080").
			Post("/crawl").
			JSON(map[string]string{"url": "http://example.com/about"}).
			Reply(200)

		err := s.Crawl(context.Background(), "http://example.com/home")

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		p, err := s.pageService.Get(page.NewPageID("http://example.com/home"))

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if p.URL != "http://example.com/home" {
			t.Errorf("expected http://example.com/home but got %s", p.URL)
		}

		if len(p.Versions) != 1 {
			t.Errorf("expected 1 version but got %d", len(p.Versions))
		}

		if p.Versions[0].Hash != page.NewVersionHash(testFile) {
			t.Errorf("expected %s but got %s", page.NewVersionHash(testFile), p.Versions[0].Hash)
		}

		data, err := s.storageService.Read(page.NewVersionHash(testFile))

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if string(data) != string(testFile) {
			t.Errorf("expected %s but got %s", testFile, data)
		}

		if !gock.IsDone() {
			t.Errorf("Not all expectations were met")
		}
	})
}
