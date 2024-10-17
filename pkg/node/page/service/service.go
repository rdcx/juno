package service

import "juno/pkg/node/page"

type Service struct {
	repo page.Repository
}

func New(repo page.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Iterator(fn func(*page.Page)) error {
	return s.repo.Iterator(fn)
}

func (s *Service) GetByURL(url string) (*page.Page, error) {
	return s.repo.GetPage(page.NewPageID(url))
}

func (s *Service) Get(id page.PageID) (*page.Page, error) {
	return s.repo.GetPage(id)
}

func (s *Service) Create(page *page.Page) error {
	return s.repo.CreatePage(page)
}

func (s *Service) AddVersion(pageID page.PageID, version page.Version) error {
	return s.repo.AddVersion(pageID, version)
}

func (s *Service) GetVersions(pageID page.PageID) ([]page.Version, error) {
	return s.repo.GetVersions(pageID)
}
