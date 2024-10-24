package service

import (
	"juno/pkg/api/extractor/strategy"

	"github.com/google/uuid"
)

type Service struct {
	repo strategy.Repository
}

func New(repo strategy.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID uuid.UUID, name, selector string, filters []*strategy.Filter) (*strategy.Strategy, error) {
	return nil, nil
}

func (s *Service) Get(id uuid.UUID) (*strategy.Strategy, error) {
	return nil, nil
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	return nil, nil
}
