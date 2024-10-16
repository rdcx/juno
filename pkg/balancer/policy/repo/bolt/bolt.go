package bolt

import (
	"encoding/json"
	"fmt"

	"juno/pkg/balancer/policy"

	bolt "go.etcd.io/bbolt"
)

type Repository struct {
	db *bolt.DB
}

// New initializes a new BoltDB repository.
func New(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create a bucket for policies if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("policies"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// Get retrieves a CrawlPolicy by hostname from the BoltDB store.
func (r *Repository) Get(hostname string) (*policy.CrawlPolicy, error) {
	var p *policy.CrawlPolicy

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("policies"))
		v := b.Get([]byte(hostname))

		if v == nil {
			return policy.ErrPolicyNotFound
		}

		// Deserialize the policy
		err := json.Unmarshal(v, &p)
		if err != nil {
			return fmt.Errorf("failed to unmarshal policy: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return p, nil
}

// Set stores a HostnamePolicy by hostname in the BoltDB store.
func (r *Repository) Set(hostname string, p *policy.CrawlPolicy) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("policies"))

		// Serialize the policy
		data, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal policy: %w", err)
		}

		// Store the serialized policy
		return b.Put([]byte(hostname), data)
	})
}
