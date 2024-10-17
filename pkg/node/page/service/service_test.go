package service

import (
	"testing"
	"time"

	"juno/pkg/node/page"
	"juno/pkg/node/page/repo/mem"
)

func setupTestService() *Service {
	repo := mem.New() // Use in-memory repo for testing
	return New(repo)
}

func TestService_Iterator(t *testing.T) {
	service := setupTestService()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := service.Create(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Create a test page using NewPageID
	anotherPage := &page.Page{
		ID:  page.NewPageID("https://example.com/another"),
		URL: "https://example.com/another",
	}
	err = service.Create(anotherPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Iterate over pages
	pages := []*page.Page{}
	err = service.Iterator(func(p *page.Page) {
		pages = append(pages, p)
	})
	if err != nil {
		t.Fatalf("failed to iterate over pages: %v", err)
	}

	// Verify that both pages were iterated over
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestService_CreatePage(t *testing.T) {
	service := setupTestService()

	// Create a test page using NewPageID
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}

	// Test creating a page
	err := service.Create(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Verify that the page is saved by getting it by URL
	retrievedPage, err := service.GetByURL("https://example.com")
	if err != nil {
		t.Fatalf("failed to retrieve page: %v", err)
	}
	if retrievedPage.ID != testPage.ID || retrievedPage.URL != testPage.URL {
		t.Errorf("expected page %v, got %v", testPage, retrievedPage)
	}
}

func TestService_AddVersion(t *testing.T) {
	service := setupTestService()

	// Create a test page
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := service.Create(testPage)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Create a test version using NewVersionHash
	testVersion := page.Version{
		Hash:      page.NewVersionHash([]byte("test version data")),
		CreatedAt: time.Now(),
	}

	// Add a version to the page
	err = service.AddVersion(testPage.ID, testVersion)
	if err != nil {
		t.Fatalf("failed to add version: %v", err)
	}

	// Verify that the version was added
	retrievedVersions, err := service.GetVersions(testPage.ID)
	if err != nil {
		t.Fatalf("failed to retrieve versions: %v", err)
	}
	if len(retrievedVersions) != 1 {
		t.Errorf("expected 1 version, got %d", len(retrievedVersions))
	}
	if retrievedVersions[0].Hash != testVersion.Hash {
		t.Errorf("expected version hash %s, got %s", testVersion.Hash, retrievedVersions[0].Hash)
	}
}

func TestService_GetVersions(t *testing.T) {
	service := setupTestService()

	// Create a test page
	testPage := &page.Page{
		ID:  page.NewPageID("https://example.com"),
		URL: "https://example.com",
	}
	err := service.Create(testPage)
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
		err := service.AddVersion(testPage.ID, version)
		if err != nil {
			t.Fatalf("failed to add version: %v", err)
		}
	}

	// Retrieve and verify versions
	retrievedVersions, err := service.GetVersions(testPage.ID)
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

func TestService_GetPageNotFound(t *testing.T) {
	service := setupTestService()

	// Try to get a non-existent page
	_, err := service.Get(page.NewPageID("https://nonexistent.com"))
	if err == nil {
		t.Errorf("expected error when retrieving non-existent page, got nil")
	}
}
