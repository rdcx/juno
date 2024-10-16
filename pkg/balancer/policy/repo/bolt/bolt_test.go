package bolt

import (
	"os"
	"testing"
	"time"

	"encoding/json"
	"juno/pkg/balancer/policy"
)

func setupTestRepo(t *testing.T) (*Repository, func()) {
	// Create a temporary BoltDB file for testing
	dbPath := "test_policy.db"

	repo, err := New(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		repo.db.Close()
		os.Remove(dbPath)
	}

	return repo, cleanup
}

func TestRepository_SetAndGet(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test policy
	hostname := "example.com"
	expectedPolicy := &policy.CrawlPolicy{
		Hostname:      hostname,
		LastCrawled:   time.Now(),
		CrawlInterval: 0,
	}

	// Set the policy
	err := repo.Set(hostname, expectedPolicy)
	if err != nil {
		t.Fatalf("failed to set policy: %v", err)
	}

	// Get the policy
	actualPolicy, err := repo.Get(hostname)
	if err != nil {
		t.Fatalf("failed to get policy: %v", err)
	}

	// Compare the two policies
	expectedJSON, _ := json.Marshal(expectedPolicy)
	actualJSON, _ := json.Marshal(actualPolicy)
	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("expected %s, got %s", string(expectedJSON), string(actualJSON))
	}
}

func TestRepository_GetNonExistent(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Try to get a policy that doesn't exist
	_, err := repo.Get("nonexistent.com")
	if err == nil {
		t.Errorf("expected error when getting nonexistent policy, got nil")
	}

	if err != policy.ErrPolicyNotFound {
		t.Errorf("expected policy.ErrPolicyNotFound, got %v", err)
	}
}
