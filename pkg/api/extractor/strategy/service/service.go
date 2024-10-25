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

func (s *Service) Create(userID, selectorID uuid.UUID, name string, t strategy.StrategyType) (*strategy.Strategy, error) {
	sel := &strategy.Strategy{
		ID:         uuid.New(),
		UserID:     userID,
		SelectorID: selectorID,
		Name:       name,
		Type:       t,
	}

	err := sel.Validate()

	if err != nil {
		return nil, err
	}

	err = s.repo.Create(sel)

	if err != nil {
		return nil, err
	}

	return sel, nil
}

func (s *Service) Get(id uuid.UUID) (*strategy.Strategy, error) {
	return s.repo.Get(id)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	return s.repo.ListByUserID(userID)
}

func (s *Service) ListBySelectorID(selectorID uuid.UUID) ([]*strategy.Strategy, error) {
	return s.repo.ListBySelectorID(selectorID)
}
