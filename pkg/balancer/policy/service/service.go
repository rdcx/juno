package service

import (
	"juno/pkg/balancer/policy"
	"time"
)

type Service struct {
	repo policy.Repository
}

func New(repo policy.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CanCrawl(p *policy.CrawlPolicy) bool {
	return time.Since(p.LastCrawled) > p.CrawlInterval
}

func (s *Service) Get(hostname string) (*policy.CrawlPolicy, error) {
	return s.repo.Get(hostname)
}

func (s *Service) RecordCrawl(p *policy.CrawlPolicy) error {
	p.LastCrawled = time.Now()
	p.TimesCrawled++

	return s.repo.Set(p.Hostname, p)
}

func (s *Service) Set(hostname string, p *policy.CrawlPolicy) error {
	return s.repo.Set(hostname, p)
}
