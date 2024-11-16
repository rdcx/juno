package service

import (
	"juno/pkg/node/page"
	"testing"
)

type mockPageService struct{}

func (s *mockPageService) Iterator(fn func(*page.Page)) error {
	return nil
}

func (s *mockPageService) GetByURL(url string) (*page.Page, error) {
	return nil, nil
}

func (s *mockPageService) Get(id page.PageID) (*page.Page, error) {
	return nil, nil
}

func (s *mockPageService) Create(page *page.Page) error {
	return nil
}

func (s *mockPageService) AddVersion(pageID page.PageID, version page.Version) error {
	return nil
}

func (s *mockPageService) GetVersions(pageID page.PageID) ([]page.Version, error) {
	return nil, nil
}

func (s *mockPageService) Count() (int, error) {
	return 10, nil
}

func TestGetInfo(t *testing.T) {
	s := New(&mockPageService{})
	i, err := s.GetInfo()
	if err != nil {
		t.Fatal(err)
	}
	if i.PageCount != 10 {
		t.Fatalf("expected 10, got %d", i.PageCount)
	}
}
