package page

import (
	"encoding/json"
	"fmt"

	"juno/pkg/node/page"

	bolt "go.etcd.io/bbolt"
)

type Repository struct {
	db *bolt.DB
}

// New initializes a new BoltDB-based Repository.
func New(dbPath string) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create a bucket for pages if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pages"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// CreatePage adds a new page to the BoltDB store.
func (r *Repository) CreatePage(p *page.Page) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pages"))

		// Serialize the Page struct to JSON
		data, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal page: %w", err)
		}

		// Store the page in BoltDB with ID as the key
		return b.Put(p.ID[:], data)
	})
}

// GetPage retrieves a page by its ID from the BoltDB store.
func (r *Repository) GetPage(id page.PageID) (*page.Page, error) {
	var p page.Page

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pages"))

		// Fetch the page by ID
		data := b.Get(id[:])
		if data == nil {
			return page.ErrPageNotFound
		}

		// Deserialize JSON into a Page struct
		if err := json.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// AddVersion adds a new version to an existing page and updates the current version.
func (r *Repository) AddVersion(pageID page.PageID, version page.Version) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pages"))

		// Fetch the page by ID
		data := b.Get(pageID[:])
		if data == nil {
			return page.ErrPageNotFound
		}

		var p page.Page
		// Deserialize the page
		if err := json.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}

		// Add the new version to the list of versions
		p.Versions = append(p.Versions, version)

		// Serialize and save the updated page
		updatedData, err := json.Marshal(&p)
		if err != nil {
			return fmt.Errorf("failed to marshal updated page: %w", err)
		}

		return b.Put(pageID[:], updatedData)
	})
}

// GetVersions retrieves all versions of a page by its ID.
func (r *Repository) GetVersions(pageID page.PageID) ([]page.Version, error) {
	var versions []page.Version

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pages"))

		// Fetch the page by ID
		data := b.Get(pageID[:])
		if data == nil {
			return page.ErrPageNotFound
		}

		var p page.Page
		// Deserialize the page
		if err := json.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("failed to unmarshal page: %w", err)
		}

		versions = p.Versions
		return nil
	})

	if err != nil {
		return nil, err
	}

	return versions, nil
}

// Close closes the BoltDB connection.
func (r *Repository) Close() error {
	return r.db.Close()
}
