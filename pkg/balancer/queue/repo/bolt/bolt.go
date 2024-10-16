package bolt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type Repository struct {
	db *bolt.DB
}

// NewURLQueueRepository initializes the BoltDB store for the queue.
func New(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create the bucket for the queue if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("url_queue"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// Push adds a URL to the queue by appending it to the end.
func (r *Repository) Push(url string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("url_queue"))

		// Get the last sequence number
		seq, err := b.NextSequence()
		if err != nil {
			return err
		}

		// Use the sequence number as the implicit key
		return b.Put(itob(seq), []byte(url))
	})
}

// Pop removes and returns the first URL in the queue.
func (r *Repository) Pop() (string, error) {
	var url string
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("url_queue"))

		// Start from the first item in the bucket
		cursor := b.Cursor()
		key, value := cursor.First()

		if key == nil {
			return fmt.Errorf("queue is empty")
		}

		// Store the first value (URL)
		url = string(value)

		// Delete the item from the queue
		return b.Delete(key)
	})
	if err != nil {
		return "", err
	}
	return url, nil
}

// itob converts an integer to a byte slice (for BoltDB keys).
func itob(v uint64) []byte {
	return []byte(fmt.Sprintf("%d", v))
}
