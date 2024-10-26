package service

import (
	"encoding/json"
	"fmt"
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/strategy"
	"juno/pkg/api/ranag"
	"juno/pkg/ranag/client"
	"os"

	"github.com/google/uuid"
)

type Service struct {
	jobRepo         job.Repository
	strategyService strategy.Service
	ranagService    ranag.Service
}

func New(jobRepo job.Repository, strategyService strategy.Service, ranagService ranag.Service) *Service {

	return &Service{
		jobRepo:         jobRepo,
		strategyService: strategyService,
		ranagService:    ranagService,
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

	ranges, err := s.ranagService.GroupByRange()

	if err != nil {
		return err
	}

	if len(ranges) == 0 {
		return fmt.Errorf("no ranges found")
	}

	var data []map[string]interface{}

	for _, r := range ranges {
		for _, ran := range r {
			client := client.New(ran.Address)

			res, err := client.SendRangeAggregationRequest(
				0,
				10000,
				strat.Selectors,
				strat.Fields,
				strat.Filters,
			)

			if err != nil {
				return err
			}

			data = append(data, res.Aggregations...)
		}
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	os.WriteFile("data.json", jsonData, 0644)

	return nil
}

func (s *Service) ProcessPending() error {
	jobs, err := s.jobRepo.ListByStatus(job.PendingStatus)

	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		fmt.Println("no pending jobs")
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
