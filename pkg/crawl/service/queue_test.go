package service

import "testing"

func TestQueue(t *testing.T) {
	t.Run("can add and pop urls", func(t *testing.T) {
		q := NewQueue()
		q.Push("http://example.com")
		q.Push("http://example.com/2")

		if url := q.Pop(); url != "http://example.com" {
			t.Errorf("expected http://example.com, got %s", url)
		}

		if url := q.Pop(); url != "http://example.com/2" {
			t.Errorf("expected http://example.com/2, got %s", url)
		}
	})
}
