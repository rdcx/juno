package service

import (
	"context"
	"juno/pkg/node/page"
	"juno/pkg/shard"
)

const CRAWL_TIMEOUT = time.Second * 10

type Service struct {
	balancers [shard.SHARDS][]string
	pageService page.Service
}

func New() *Service {
	return &Service{}
}

func (s *Service) Crawl(ctx context.Context, urlStr string) error {
	body, status, finalURL, err := FetchPage(context.WithTimeout(ctx, CRAWL_TIMEOUT), urlStr)

	if err != nil {
		return err
	}

	if status != 200 {
		return crawl.ErrNon200
	}

	links, err := s.domService.ExtractLinks(body)

	if err != nil {
		return err
	}






