package service

import (
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/ranag"

	"github.com/google/uuid"
)

type Service struct {
	jobRepo         job.Repository
	strategyService strategy.Service
	ranagService    ranag.Service

	ranags map[[2]int][]string
}

func New(jobRepo job.Repository, strategyService strategy.Service) *Service {

	ranges := 10000

	fakeRanags := make(map[[2]int][]string)

	for i := 0; i < ranges; i++ {
		for j := 0; j < ranges; j++ {
			fakeRanags[[2]int{i, j}] = []string{"localhost:9292"}
		}
	}

	return &Service{
		jobRepo:         jobRepo,
		strategyService: strategyService,
		ranags:          fakeRanags,
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

func (s *Service) process(j *job.Job) error {
	strat, err := s.strategyService.Get(j.StrategyID)

	if err != nil {
		return err
	}

}

func (s *Service) ProcessPending() error {
	jobs, err := s.jobRepo.ListByStatus(job.PendingStatus)

	if err != nil {
		return err
	}

	for _, j := range jobs {

		j.Status = job.RunningStatus

		err := s.jobRepo.Update(j)

		if err != nil {
			return err
		}

		err = s.process(j)

		if err != nil {
			j.Status = job.FailedStatus
		} else {
			j.Status = job.CompletedStatus
		}

		err = s.jobRepo.Update(j)

		if err != nil {
			return err
		}
	}

	return nil
}
