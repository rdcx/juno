package service

import (
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/strategy"

	"github.com/google/uuid"
)

type Service struct {
	jobRepo         job.Repository
	strategyService strategy.Service
}

func New(jobRepo job.Repository, strategyService strategy.Service) *Service {
	return &Service{
		jobRepo:         jobRepo,
		strategyService: strategyService,
	}
}

func (s *Service) Get(id uuid.UUID) (*job.Job, error) {
	return s.jobRepo.Get(id)
}

func (s *Service) Create(userID, strategyID uuid.UUID) (*job.Job, error) {

	e, err := s.strategyService.Get(strategyID)

	if err != nil {
		return nil, err
	}

	q := &job.Job{
		ID:         uuid.New(),
		UserID:     userID,
		StrategyID: e.ID,
		Status:     job.PendingStatus,
	}

	err = s.jobRepo.Create(q)

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (s *Service) Update(q *job.Job) error {
	_, err := s.jobRepo.Get(q.ID)

	if err != nil {
		return err
	}
	return s.jobRepo.Update(q)
}

func (s *Service) ListByUserID(userID uuid.UUID) ([]*job.Job, error) {
	return s.jobRepo.ListByUserID(userID)
}
