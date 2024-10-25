package service

import (
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/job/repo/mem"
	"juno/pkg/api/extractor/strategy"
	"testing"

	"github.com/google/uuid"
)

type mockStrategyService struct {
	returnStrategy *strategy.Strategy
	returnError    error
}

func (m *mockStrategyService) Get(id uuid.UUID) (*strategy.Strategy, error) {
	return m.returnStrategy, m.returnError
}

func (m *mockStrategyService) Create(userID uuid.UUID, name string) (*strategy.Strategy, error) {
	return m.returnStrategy, m.returnError
}

func (m *mockStrategyService) ListByUserID(userID uuid.UUID) ([]*strategy.Strategy, error) {
	return nil, nil
}

func (m *mockStrategyService) Update(e *strategy.Strategy) error {
	return nil
}

func (m *mockStrategyService) Delete(id uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) AddField(id uuid.UUID, fieldID uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) RemoveField(id uuid.UUID, fieldID uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) AddFilter(id uuid.UUID, filterID uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) RemoveFilter(id uuid.UUID, filterID uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) AddSelector(id uuid.UUID, selectorID uuid.UUID) error {
	return nil
}

func (m *mockStrategyService) RemoveSelector(id uuid.UUID, selectorID uuid.UUID) error {
	return nil
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		strategyID := uuid.New()
		service := New(repo, &mockStrategyService{
			returnStrategy: &strategy.Strategy{
				ID: strategyID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, strategyID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if j.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, j.UserID)
		}

		if j.StrategyID != strategyID {
			t.Errorf("Expected %s, got %s", strategyID, j.StrategyID)
		}

		if j.Status != job.PendingStatus {
			t.Errorf("Expected %s, got %s", job.PendingStatus, j.Status)
		}

		if j.ID == uuid.Nil {
			t.Errorf("Expected non-zero UUID, got zero")
		}

		check, err := repo.Get(j.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}

		if check.StrategyID != j.StrategyID {
			t.Errorf("Expected %s, got %s", j.StrategyID, check.StrategyID)
		}
	})

	t.Run("strategy not found", func(t *testing.T) {
		repo := mem.New()
		strategyID := uuid.New()
		service := New(repo, &mockStrategyService{
			returnError: strategy.ErrNotFound,
		})
		userID := uuid.New()
		_, err := service.Create(userID, strategyID)

		if err != strategy.ErrNotFound {
			t.Errorf("Expected %v, got %v", strategy.ErrNotFound, err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		strategyID := uuid.New()
		service := New(repo, &mockStrategyService{
			returnStrategy: &strategy.Strategy{
				ID: strategyID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, strategyID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := service.Get(j.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, check.ID)
		}

		if check.UserID != j.UserID {
			t.Errorf("Expected %s, got %s", j.UserID, check.UserID)
		}

		if check.Status != j.Status {
			t.Errorf("Expected %s, got %s", j.Status, check.Status)
		}

		if check.StrategyID != j.StrategyID {
			t.Errorf("Expected %s, got %s", j.StrategyID, check.StrategyID)
		}
	})

	t.Run("job not found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo, &mockStrategyService{})
		_, err := service.Get(uuid.New())

		if err != job.ErrNotFound {
			t.Errorf("Expected %v, got %v", job.ErrNotFound, err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		strategyID := uuid.New()
		service := New(repo, &mockStrategyService{
			returnStrategy: &strategy.Strategy{
				ID: strategyID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, strategyID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		list, err := service.ListByUserID(userID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(list) != 1 {
			t.Errorf("Expected 1, got %d", len(list))
		}

		if list[0].ID != j.ID {
			t.Errorf("Expected %s, got %s", j.ID, list[0].ID)
		}
	})

	t.Run("no jobs found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo, &mockStrategyService{})
		list, err := service.ListByUserID(uuid.New())

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(list) != 0 {
			t.Errorf("Expected 0, got %d", len(list))
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		strategyID := uuid.New()
		service := New(repo, &mockStrategyService{
			returnStrategy: &strategy.Strategy{
				ID: strategyID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, strategyID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		j.Status = job.CompletedStatus

		err = service.Update(j)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		check, err := repo.Get(j.ID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if check.Status != job.CompletedStatus {
			t.Errorf("Expected %s, got %s", job.CompletedStatus, check.Status)
		}
	})

	t.Run("job not found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo, &mockStrategyService{})
		err := service.Update(&job.Job{})

		if err != job.ErrNotFound {
			t.Errorf("Expected %v, got %v", job.ErrNotFound, err)
		}
	})
}

func TestProcessPending(t *testing.T) {
	t.Run("sets job status to running", func(t *testing.T) {

}