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

	// Channel to collect results from goroutines
	results := make(chan []map[string]interface{}, len(ranges))
	errors := make(chan error, len(ranges))

	// Semaphore channel to limit the number of goroutines
	sem := make(chan struct{}, 10) // 10 goroutines

	// Launch goroutines
	for rval, r := range ranges {
		for _, ran := range r {
			sem <- struct{}{} // acquire a token
			go func(rval [2]int, ran *ranag.Ranag) {
				defer func() { <-sem }() // release token

				client := client.New(ran.Address)
				res, err := client.SendRangeAggregationRequest(
					rval[0],
					rval[1],
					strat.Selectors,
					strat.Fields,
					strat.Filters,
				)
				if err != nil {
					errors <- err
					return
				}

				results <- res.Aggregations
			}(rval, ran)
		}
	}

	// Wait for all goroutines to complete
	go func() {
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{} // fill up semaphore channel to ensure all goroutines are done
		}
		close(results)
		close(errors)
	}()

	// Collect results
	var data []map[string]interface{}
	for res := range results {
		data = append(data, res...)
	}

	// Check if any error occurred
	if len(errors) > 0 {
		return <-errors
	}

	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile("data.json", jsonData, 0644); err != nil {
		return err
	}

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
			fmt.Println(err)
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
