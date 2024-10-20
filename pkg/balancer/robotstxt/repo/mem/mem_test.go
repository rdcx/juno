package mem

import "testing"

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := New()

		repo.Set("http://example.com", "User-agent: *\nDisallow: /private")

		hit, err := repo.Get("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if hit != "User-agent: *\nDisallow: /private" {
			t.Errorf("expected hit to be User-agent: *\nDisallow: /private, got %s", hit)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := New()

		err := repo.Set("http://example.com", "User-agent: *\nDisallow: /private")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		hit, err := repo.Get("http://example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if hit != "User-agent: *\nDisallow: /private" {
			t.Errorf("expected hit to be User-agent: *\nDisallow: /private, got %s", hit)
		}
	})
}
