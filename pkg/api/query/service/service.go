package service

import (
	"juno/pkg/api/query"

	"github.com/google/uuid"
)

type Service struct {
	queryRepo query.Repository
}

func New(queryRepo query.Repository) *Service {
	return &Service{
		queryRepo: queryRepo,
	}
}

func (s *Service) Get(id uuid.UUID) (*query.Query, error) {
	return s.queryRepo.Get(id)
}

func (s *Service) Create(userID uuid.UUID, qt query.QueryType, bq *query.BasicQuery) (*query.Query, error) {

	q := &query.Query{
		ID:                uuid.New(),
		UserID:            userID,
		Status:            query.PendingStatus,
		QueryType:         qt,
		BasicQueryVersion: "v1",
		BasicQuery:        bq,
	}
	err := s.queryRepo.Create(q)

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (s *Service) Update(q *query.Query) error {
	return s.queryRepo.Update(q)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*query.Query, error) {
	return s.queryRepo.ListByUserID(userID)
}
