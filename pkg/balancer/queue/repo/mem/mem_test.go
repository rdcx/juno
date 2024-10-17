package mem

import "testing"

func TestExists(t *testing.T) {
	t.Run("should return false when url does not exist", func(t *testing.T) {
		// Given
		repo := New()
		url := "http://example.com"

		// When
		exists, err := repo.Exists(url)

		// Then
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if exists {
			t.Errorf("expected false but got true")
		}
	})

	t.Run("should return true when url exists", func(t *testing.T) {
		// Given
		repo := New()
		url := "http://example.com"
		repo.urls = []string{url}

		// When
		exists, err := repo.Exists(url)

		// Then
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if !exists {
			t.Errorf("expected true but got false")
		}
	})
}

func TestPush(t *testing.T) {
	t.Run("should push url", func(t *testing.T) {
		// Given
		repo := New()
		url := "http://example.com"

		// When
		err := repo.Push(url)

		// Then
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if len(repo.urls) != 1 {
			t.Errorf("expected 1 url but got %d", len(repo.urls))
		}

		if repo.urls[0] != url {
			t.Errorf("expected %s but got %s", url, repo.urls[0])
		}
	})
}

func TestPop(t *testing.T) {
	t.Run("should pop url", func(t *testing.T) {
		// Given
		repo := New()
		url := "http://example.com"
		repo.urls = []string{url}

		// When
		popped, err := repo.Pop()

		// Then
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if len(repo.urls) != 0 {
			t.Errorf("expected 0 urls but got %d", len(repo.urls))
		}

		if popped != url {
			t.Errorf("expected %s but got %s", url, popped)
		}
	})

	t.Run("should return error when no urls", func(t *testing.T) {
		// Given
		repo := New()

		// When
		_, err := repo.Pop()

		// Then
		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})
}
