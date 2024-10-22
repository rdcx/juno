package service

import (
	"juno/pkg/api/extraction/extractor"

	"github.com/google/uuid"
)

type Service struct {
	repo extractor.Repository
}

func New(repo extractor.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID uuid.UUID, name, selector string, filters []*extractor.Filter) (*extractor.Extractor, error) {
	return nil, nil
}

func (s *Service) Get(id uuid.UUID) (*extractor.Extractor, error) {
	return nil, nil
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*extractor.Extractor, error) {
	return nil, nil
}
