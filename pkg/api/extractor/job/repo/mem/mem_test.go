package mem

import (
	"juno/pkg/api/extractor/job"
	"testing"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	repo := New()

	q := &job.Job{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Status:     job.PendingStatus,
		StrategyID: uuid.New(),
	}

	err := repo.Create(q)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check, err := repo.Get(q.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, check.ID)
	}

	if check.UserID != q.UserID {
		t.Errorf("Expected %s, got %s", q.UserID, check.UserID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	if check.StrategyID != q.StrategyID {
		t.Errorf("Expected %s, got %s", q.StrategyID, check.StrategyID)
	}

	if check.Status != q.Status {
		t.Errorf("Expected %s, got %s", q.Status, check.Status)
	}

	list, err := repo.ListByUserID(q.UserID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}

	if list[0].ID != q.ID {
		t.Errorf("Expected %s, got %s", q.ID, list[0].ID)
	}

	var copied job.Job

	copied.ID = q.ID
	copied.UserID = q.UserID
	copied.Status = q.Status
	copied.StrategyID = q.StrategyID

	copied.Status = job.CompletedStatus

	err = repo.Update(&copied)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	check, err = repo.Get(q.ID)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if check.Status != job.CompletedStatus {
		t.Errorf("Expected %s, got %s", job.CompletedStatus, check.Status)
	}

	list, err = repo.ListByStatus(job.CompletedStatus)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1, got %d", len(list))
	}
}
