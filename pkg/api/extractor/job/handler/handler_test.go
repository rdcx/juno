package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"juno/pkg/api/auth"
	"juno/pkg/api/extractor/job"
	"juno/pkg/api/extractor/job/dto"
	"juno/pkg/api/extractor/job/policy"
	"juno/pkg/api/user"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockJobService struct {
	withError  error
	withUserID uuid.UUID
}

func (m *mockJobService) Create(userID, strategyID uuid.UUID) (*job.Job, error) {

	if m.withError != nil {
		return nil, m.withError
	}

	return &job.Job{
		ID:         uuid.New(),
		UserID:     userID,
		StrategyID: strategyID,
		Status:     job.PendingStatus,
	}, nil
}

func (m *mockJobService) Get(id uuid.UUID) (*job.Job, error) {
	if m.withError != nil {
		return nil, m.withError
	}
	return &job.Job{
		ID:         id,
		UserID:     m.withUserID,
		StrategyID: uuid.New(),
		Status:     job.PendingStatus,
	}, nil
}

func (m *mockJobService) ListByUserID(userID uuid.UUID) ([]*job.Job, error) {
	if m.withError != nil {
		return nil, m.withError
	}
	return []*job.Job{
		{
			ID:         uuid.New(),
			UserID:     userID,
			StrategyID: uuid.New(),
			Status:     job.PendingStatus,
		},
	}, nil
}

func (m *mockJobService) Update(q *job.Job) error {
	if m.withError != nil {
		return m.withError
	}
	return nil
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		h := New(&mockJobService{}, policy.New())

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		req := dto.CreateJobRequest{
			StrategyID: uuid.New().String(),
		}

		encoded, err := json.Marshal(req)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		userID := uuid.New()
		c.Request = httptest.NewRequest("POST", "/jobs", bytes.NewBuffer(encoded)).WithContext(
			auth.WithUser(context.Background(), &user.User{
				ID: userID,
			}),
		)

		h.Create(c)

		if w.Code != 201 {
			t.Errorf("Expected 201, got %d", w.Code)
		}

		var res dto.CreateJobResponse

		b, _ := io.ReadAll(w.Body)

		err = json.Unmarshal(b, &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Job.ID == "" {
			t.Errorf("Expected not empty, got empty")
		}

		if res.Job.Status != string(job.PendingStatus) {
			t.Errorf("Expected %s, got %s", job.PendingStatus, res.Job.Status)
		}

		if res.Job.UserID != userID.String() {
			t.Errorf("Expected %s, got %s", userID.String(), res.Job.UserID)
		}

		if res.Job.StrategyID != req.StrategyID {
			t.Errorf("Expected %s, got %s", req.StrategyID, res.Job.StrategyID)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		h := New(&mockJobService{
			withUserID: userID,
		}, policy.New())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jobID := uuid.New()

		c.Request = httptest.NewRequest("GET", "/jobs/"+jobID.String(), nil).WithContext(
			auth.WithUser(context.Background(), &user.User{
				ID: userID,
			}),
		)
		c.Params = gin.Params{{Key: "id", Value: jobID.String()}}

		h.Get(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var res dto.GetJobResponse
		b, _ := io.ReadAll(w.Body)
		err := json.Unmarshal(b, &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if res.Status != dto.SUCCESS {
			t.Fatalf("Expected %s, got %s", dto.SUCCESS, res.Status)
		}

		if res.Job.ID == "" {
			t.Errorf("Expected not empty, got empty")
		}

		if res.Job.ID != jobID.String() {
			t.Errorf("Expected job ID %s, got %s", jobID.String(), res.Job.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := New(&mockJobService{withError: job.ErrNotFound}, policy.New())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jobID := uuid.New()
		userID := uuid.New()

		c.Request = httptest.NewRequest("GET", "/jobs/"+jobID.String(), nil).WithContext(
			auth.WithUser(context.Background(), &user.User{
				ID: userID,
			}),
		)
		c.Params = gin.Params{{Key: "id", Value: jobID.String()}}

		h.Get(c)

		if w.Code != 404 {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := New(&mockJobService{}, policy.New())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Request = httptest.NewRequest("GET", "/jobs", nil).WithContext(
			auth.WithUser(context.Background(), &user.User{
				ID: userID,
			}),
		)

		h.List(c)

		if w.Code != 200 {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var res dto.ListJobsResponse
		b, _ := io.ReadAll(w.Body)
		err := json.Unmarshal(b, &res)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if len(res.Jobs) == 0 {
			t.Errorf("Expected jobs, got empty list")
		}

		if res.Jobs[0].UserID != userID.String() {
			t.Errorf("Expected user ID %s, got %s", userID.String(), res.Jobs[0].UserID)
		}
	})
}
