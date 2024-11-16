package page

import (
	"os"
	"testing"
	"time"

	"juno/pkg/node/page"
)

// setupTestRepo sets up a BoltDB repository for testing and returns a cleanup function.
func setupTestRepo(t *testing.T) (*Repository, func()) {
	// Create a temporary BoltDB file for testing
	dbPath := "test_pages.db"
	repo, err := New(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Cleanup function to close the database and remove the file
	cleanup := func() {
		repo.Close()
		os.Remove(dbPath)
	}

	return repo, cleanup
}

func TestRepository_Iterator(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := repo.CreatePage(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Create a test page using NewPageID
	anotherPage := &page.Page{
		ID:  page.NewPageID("https://example.com/another"),
		URL: "https://example.com/another",
	}
	err = repo.CreatePage(anotherPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Iterate over pages
	pages := []*page.Page{}
	err = repo.Iterator(
		func(p *page.Page) {
			pages = append(pages, p)
		},
	)
	if err != nil {
		t.Fatalf("failed to iterate over pages: %v", err)
	}

	// Verify that both pages were iterated over
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestRepository_CreatePage(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}

	// Test creating a page
	err := repo.CreatePage(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Verify that the page is saved
	retrievedPage, err := repo.GetPage(testPage.ID)
	if err != nil {
		t.Fatalf("failed to retrieve page: %v", err)
	}
	if retrievedPage.ID != testPage.ID || retrievedPage.URL != testPage.URL {
		t.Errorf("expected page %v, got %v", testPage, retrievedPage)
	}
}

func TestRepository_AddVersion(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := repo.CreatePage(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Create a test version using NewVersionHash
	testVersion := page.Version{
		Hash:      page.NewVersionHash([]byte("test version data")),
		CreatedAt: time.Now(),
	}

	// Add a version to the page
	err = repo.AddVersion(testPage.ID, testVersion)
	if err != nil {
		t.Fatalf("failed to add version: %v", err)
	}

	// Verify that the version was added
	retrievedPage, err := repo.GetPage(testPage.ID)
	if err != nil {
		t.Fatalf("failed to retrieve page: %v", err)
	}
	if len(retrievedPage.Versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(retrievedPage.Versions))
	}
	if retrievedPage.Versions[0].Hash != testVersion.Hash {
		t.Errorf("expected version hash %s, got %s", testVersion.Hash, retrievedPage.Versions[0].Hash)
	}
}

func TestRepository_GetVersions(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := repo.CreatePage(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Add some versions to the page using NewVersionHash
	versions := []page.Version{
		{
			Hash:      page.NewVersionHash([]byte("version 1")),
			CreatedAt: time.Now(),
		},
		{
			Hash:      page.NewVersionHash([]byte("version 2")),
			CreatedAt: time.Now(),
		},
	}

	for _, version := range versions {
		err := repo.AddVersion(testPage.ID, version)
		if err != nil {
			t.Fatalf("failed to add version: %v", err)
		}
	}

	// Retrieve and verify versions
	retrievedVersions, err := repo.GetVersions(testPage.ID)
	if err != nil {
		t.Fatalf("failed to retrieve versions: %v", err)
	}
	if len(retrievedVersions) != len(versions) {
		t.Errorf("expected %d versions, got %d", len(versions), len(retrievedVersions))
	}

	for i, version := range versions {
		if retrievedVersions[i].Hash != version.Hash {
			t.Errorf("expected version hash %s, got %s", version.Hash, retrievedVersions[i].Hash)
		}
	}
}

func TestRepository_GetPageNotFound(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Try to get a non-existent page
	_, err := repo.GetPage(page.NewPageID("https://nonexistent.com"))
	if err == nil {
		t.Errorf("expected error when retrieving non-existent page, got nil")
	}
}

func TestRepository_Count(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create some test pages
	testPages := []*page.Page{
		{
			ID:  page.NewPageID("https://example.com"),
			URL: "https://example.com",
		},
		{
			ID:  page.NewPageID("https://example.com/another"),
			URL: "https://example.com/another",
		},
	}

	for _, testPage := range testPages {
		err := repo.CreatePage(testPage)
		if err != nil {
			t.Fatalf("failed to create page: %v", err)
		}
	}

	// Verify the count
	count, err := repo.Count()
	if err != nil {
		t.Fatalf("failed to get count: %v", err)
	}
	if count != len(testPages) {
		t.Errorf("expected count %d, got %d", len(testPages), count)
	}
}
