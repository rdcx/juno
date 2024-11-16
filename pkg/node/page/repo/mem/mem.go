package mem

import (
	"juno/pkg/node/page"
	"sync"
	"time"
)

type Repository struct {
	mu    sync.RWMutex // Mutex to handle concurrent access
	pages map[page.PageID]*page.Page
}

// New initializes a new in-memory repository.
func New() *Repository {
	return &Repository{
		pages: make(map[page.PageID]*page.Page),
	}
}

func (r *Repository) Iterator(callback func(*page.Page)) error {

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Example iteration logic over a collection of pages
	for _, p := range r.pages {
		// Call the provided callback function for each page
		callback(p)
	}
	return nil
}

// CreatePage adds a new page to the in-memory store.
func (r *Repository) CreatePage(p *page.Page) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the page already exists
	if _, exists := r.pages[p.ID]; exists {
		return page.ErrPageAlreadyExists
	}

	r.pages[p.ID] = p
	return nil
}

// GetPage retrieves a page by its ID from the in-memory store.
func (r *Repository) GetPage(id page.PageID) (*page.Page, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.pages[id]
	if !exists {
		return nil, page.ErrPageNotFound
	}
	return p, nil
}

// AddVersion adds a new version to an existing page.
func (r *Repository) AddVersion(pageID page.PageID, version page.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the page exists
	p, exists := r.pages[pageID]
	if !exists {
		return page.ErrPageNotFound
	}

	// Add the new version to the page
	version.CreatedAt = time.Now() // Set the current time for the version
	p.Versions = append(p.Versions, version)
	return nil
}

// GetVersions retrieves all versions of a page by its ID.
func (r *Repository) GetVersions(pageID page.PageID) ([]page.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.pages[pageID]
	if !exists {
		return nil, page.ErrPageNotFound
	}

	return p.Versions, nil
}

func (r *Repository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.pages), nil
}
