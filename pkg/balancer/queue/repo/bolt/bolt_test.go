package bolt

import (
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestNewURLQueueRepository(t *testing.T) {
	// Create a temporary BoltDB file for testing
	dbPath := "test_queue_new.db"
	defer os.Remove(dbPath) // Clean up after test

	// Initialize the queue repository
	repo, err := NewURLQueueRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.db.Close()

	// Check if the bucket was created
	err = repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("url_queue"))
		if b == nil {
			t.Errorf("bucket 'url_queue' not created")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to view bucket: %v", err)
	}
}

func TestPush(t *testing.T) {
	dbPath := "test_queue_push.db"
	defer os.Remove(dbPath) // Clean up after test

	// Initialize the queue repository
	repo, err := NewURLQueueRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.db.Close()

	// Test pushing URLs to the queue
	url1 := "https://example.com"
	url2 := "https://another-example.com"

	err = repo.Push(url1)
	if err != nil {
		t.Errorf("failed to push url1: %v", err)
	}

	err = repo.Push(url2)
	if err != nil {
		t.Errorf("failed to push url2: %v", err)
	}

	// Verify URLs were pushed
	repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("url_queue"))
		if b == nil {
			t.Errorf("bucket 'url_queue' not found")
		}

		// Check first URL
		v := b.Get(itob(1))
		if string(v) != url1 {
			t.Errorf("expected %s, got %s", url1, string(v))
		}

		// Check second URL
		v = b.Get(itob(2))
		if string(v) != url2 {
			t.Errorf("expected %s, got %s", url2, string(v))
		}
		return nil
	})
}

func TestPop(t *testing.T) {
	dbPath := "test_queue_pop.db"
	defer os.Remove(dbPath) // Clean up after test

	// Initialize the queue repository
	repo, err := NewURLQueueRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.db.Close()

	// Push some URLs for testing
	url1 := "https://example.com"
	url2 := "https://another-example.com"

	repo.Push(url1)
	repo.Push(url2)

	// Pop the first URL
	poppedURL, err := repo.Pop()
	if err != nil {
		t.Fatalf("failed to pop: %v", err)
	}
	if poppedURL != url1 {
		t.Errorf("expected %s, got %s", url1, poppedURL)
	}

	// Pop the second URL
	poppedURL, err = repo.Pop()
	if err != nil {
		t.Fatalf("failed to pop: %v", err)
	}
	if poppedURL != url2 {
		t.Errorf("expected %s, got %s", url2, poppedURL)
	}

	// Try to pop from an empty queue
	_, err = repo.Pop()
	if err == nil {
		t.Error("expected error when popping from empty queue, got nil")
	}
}
