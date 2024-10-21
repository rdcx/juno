package service

import (
	"context"
	"juno/pkg/node/balancer"
	"juno/pkg/node/crawl"
	"juno/pkg/node/fetcher"
	"juno/pkg/node/html"
	"juno/pkg/node/page"
	"juno/pkg/node/storage"
	"juno/pkg/shard"
	"juno/pkg/url"
	"time"
)

const CRAWL_TIMEOUT = time.Second * 10

type Service struct {
	balancerService balancer.Service
	pageService     page.Service
	storageService  storage.Service
	htmlService     html.Service
	fetcher         fetcher.Service
}

func New(
	balancerService balancer.Service,
	pageService page.Service,
	storageService storage.Service,
	fetcher fetcher.Service,
	htmlService html.Service,
) *Service {
	return &Service{
		balancerService: balancerService,
		pageService:     pageService,
		storageService:  storageService,
		fetcher:         fetcher,
		htmlService:     htmlService,
	}
}

func (s *Service) Crawl(ctx context.Context, urlStr string) error {
	ctx, cancel := context.WithTimeout(ctx, CRAWL_TIMEOUT)

	defer cancel()
	body, status, finalURL, err := s.fetcher.FetchPage(ctx, urlStr)

	if err != nil {
		return err
	}

	if status != 200 {
		return crawl.ErrNon200Response
	}

	p, err := s.pageService.Get(page.NewPageID(finalURL))

	if err == page.ErrPageNotFound {
		p = page.NewPage(finalURL)

		hostname, err := url.ToHostname(finalURL)

		if err != nil {
			return err
		}

		shard := shard.GetShard(hostname)
		p.Shard = shard

		err = s.pageService.Create(p)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	vHash := page.NewVersionHash(body)

	err = s.storageService.Write(
		vHash,
		body,
	)

	if err != nil {
		return err
	}

	links, err := s.htmlService.ExtractLinks(body)

	if err != nil {
		return err
	}

	var fullLinks []string

	for _, link := range links {
		full, err := url.LinkToFullURL(finalURL, link)

		if err != nil {
			continue
		}

		if !url.IsHTTPOrHTTPS(full) {
			continue
		}

		fullLinks = append(fullLinks, full)

	}

	go s.balancerService.SendBatchedLinks(fullLinks)

	err = s.pageService.AddVersion(p.ID, page.NewVersion(vHash))

	if err != nil {
		return err
	}

	return s.balancerService.ReportURLProcessed(urlStr, status)
}
