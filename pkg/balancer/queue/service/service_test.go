package service

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

type mockQueueRepo struct {
	pushedURL string
	withError error
	exists    bool
}

func (m *mockQueueRepo) Push(url string) error {
	m.pushedURL = url
	return m.withError
}

func (m *mockQueueRepo) Pop() (string, error) {
	return "", nil
}

func (m *mockQueueRepo) Exists(url string) (bool, error) {
	return m.exists, nil
}

func TestPush(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockQueueRepo{}
		service := New(logrus.New(), repo)

		err := service.Push("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if repo.pushedURL != "http://example.com" {
			t.Errorf("expected pushedURL to be http://example.com, got %s", repo.pushedURL)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockQueueRepo{withError: errors.New("repo error")}
		service := New(logrus.New(), repo)

		err := service.Push("http://example.com")

		if err.Error() != "repo error" {
			t.Errorf("expected ErrRepo, got %v", err)
		}
	})

	t.Run("exists", func(t *testing.T) {
		repo := &mockQueueRepo{exists: true}
		service := New(logrus.New(), repo)

		err := service.Push("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if repo.pushedURL != "" {
			t.Errorf("expected pushedURL to be empty, got %s", repo.pushedURL)
		}
	})
}
