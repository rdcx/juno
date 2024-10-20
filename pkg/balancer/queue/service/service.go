package service

import (
	"juno/pkg/balancer/queue"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	repo   queue.Repository
}

func New(
	logger *logrus.Logger,
	repo queue.Repository,

) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) Push(url string) error {
	exists, err := s.repo.Exists(url)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return s.repo.Push(url)
}

func (s *Service) Pop() (string, error) {
	return s.repo.Pop()
}
