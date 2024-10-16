package page

import "testing"

func TestNewPageID(t *testing.T) {
	t.Run("should return new page id", func(t *testing.T) {
		id := NewPageID("http://example.com")

		if id.String() != "f0e6a6a97042a4f1f1c87f5f7d44315b" {
			t.Errorf("expected empty f0e6a6a97042a4f1f1c87f5f7d44315b but got %s", id.String())
		}
	})
}
