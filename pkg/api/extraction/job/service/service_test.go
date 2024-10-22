package service

import (
	"juno/pkg/api/extraction/extractor"
	"juno/pkg/api/extraction/job"
	"juno/pkg/api/extraction/job/repo/mem"
	"testing"

	"github.com/google/uuid"
)

type mockExtractorService struct {
	returnExtractor *extractor.Extractor
	returnError     error
}

func (m *mockExtractorService) Get(id uuid.UUID) (*extractor.Extractor, error) {
	return m.returnExtractor, m.returnError
}

func (m *mockExtractorService) Create(userID uuid.UUID, name, selector string, filters []*extractor.Filter) (*extractor.Extractor, error) {
	return m.returnExtractor, m.returnError
}

func (m *mockExtractorService) ListByUserID(userID uuid.UUID) ([]*extractor.Extractor, error) {
	return nil, nil
}

func (m *mockExtractorService) Update(e *extractor.Extractor) error {
	return nil
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		extractorID := uuid.New()
		service := New(repo, &mockExtractorService{
			returnExtractor: &extractor.Extractor{
				ID: extractorID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, extractorID)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if j.UserID != userID {
			t.Errorf("Expected %s, got %s", userID, j.UserID)
		}

		if j.ExtractorID != extractorID {
			t.Errorf("Expected %s, got %s", extractorID, j.ExtractorID)
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

		if check.ExtractorID != j.ExtractorID {
			t.Errorf("Expected %s, got %s", j.ExtractorID, check.ExtractorID)
		}
	})

	t.Run("extractor not found", func(t *testing.T) {
		repo := mem.New()
		extractorID := uuid.New()
		service := New(repo, &mockExtractorService{
			returnError: extractor.ErrNotFound,
		})
		userID := uuid.New()
		_, err := service.Create(userID, extractorID)

		if err != extractor.ErrNotFound {
			t.Errorf("Expected %v, got %v", extractor.ErrNotFound, err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		extractorID := uuid.New()
		service := New(repo, &mockExtractorService{
			returnExtractor: &extractor.Extractor{
				ID: extractorID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, extractorID)

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

		if check.ExtractorID != j.ExtractorID {
			t.Errorf("Expected %s, got %s", j.ExtractorID, check.ExtractorID)
		}
	})

	t.Run("job not found", func(t *testing.T) {
		repo := mem.New()
		service := New(repo, &mockExtractorService{})
		_, err := service.Get(uuid.New())

		if err != job.ErrNotFound {
			t.Errorf("Expected %v, got %v", job.ErrNotFound, err)
		}
	})
}

func TestListByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		extractorID := uuid.New()
		service := New(repo, &mockExtractorService{
			returnExtractor: &extractor.Extractor{
				ID: extractorID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, extractorID)

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
		service := New(repo, &mockExtractorService{})
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
		extractorID := uuid.New()
		service := New(repo, &mockExtractorService{
			returnExtractor: &extractor.Extractor{
				ID: extractorID,
			},
		})
		userID := uuid.New()
		j, err := service.Create(userID, extractorID)

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
		service := New(repo, &mockExtractorService{})
		err := service.Update(&job.Job{})

		if err != job.ErrNotFound {
			t.Errorf("Expected %v, got %v", job.ErrNotFound, err)
		}
	})
}
