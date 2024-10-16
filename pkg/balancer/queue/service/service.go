package service

import (
	"context"
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/policy"
	"juno/pkg/balancer/queue"
	"juno/pkg/link"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger        *logrus.Logger
	repo          queue.Repository
	crawlService  crawl.Service
	policyService policy.Service
}

func New(
	logger *logrus.Logger,
	repo queue.Repository,
	crawlService crawl.Service,
	policyService policy.Service,

) *Service {
	return &Service{
		logger:        logger,
		repo:          repo,
		crawlService:  crawlService,
		policyService: policyService,
	}
}

func (s *Service) Push(url string) error {
	return s.repo.Push(url)
}

func (s *Service) Pop() (string, error) {
	return s.repo.Pop()
}

func (s *Service) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return queue.ErrProcessQueueCancelled
		default:
			url, err := s.Pop()

			if err == queue.ErrNoURLsInQueue {
				select {
				case <-ctx.Done():
					return queue.ErrProcessQueueCancelled
				case <-time.After(500 * time.Millisecond):
					// continue after sleep
				}
				continue
			}

			hostname, err := link.ToHostname(url)

			if err != nil {
				s.logger.Errorf("failed to get hostname from url: %v", err)
				continue
			}

			pol, err := s.policyService.Get(hostname)

			if err == policy.ErrPolicyNotFound {
				pol = policy.New(hostname)
			} else if err != nil {
				s.logger.Errorf("failed to get policy for url: %v", err)
				continue
			}

			if !s.policyService.CanCrawl(pol) {
				err = s.Push(url)

				if err != nil {
					s.logger.Errorf("failed to push url to queue: %v", err)
				}

				continue
			}

			err = s.crawlService.Crawl(url)

			if err == nil {
				err = s.policyService.RecordCrawl(pol)

				if err != nil {
					s.logger.Errorf("failed to set policy for url: %v", err)
				}

				continue
			}

			err = s.Push(url)

			if err != nil {
				s.logger.Errorf("failed to push url to queue: %v", err)
			}
		}
	}
}
