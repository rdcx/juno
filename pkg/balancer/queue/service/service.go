package service

import (
	"context"
	"juno/pkg/balancer/crawl"
	"juno/pkg/balancer/queue"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger       *logrus.Logger
	repo         queue.Repository
	crawlService crawl.Service
}

func New(logger *logrus.Logger, repo queue.Repository, crawlService crawl.Service) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		crawlService: crawlService,
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

			err = s.crawlService.Crawl(url)

			if err == nil {
				continue
			}

			err = s.Push(url)

			if err != nil {
				s.logger.Errorf("failed to push url to queue: %v", err)
			}
		}
	}
}
