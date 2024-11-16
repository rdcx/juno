package service

import (
	"juno/pkg/node/info"
	"juno/pkg/node/page"
)

type Service struct {
	pageService page.Service
}

func New(pageService page.Service) *Service {
	return &Service{
		pageService: pageService,
	}
}

func (s *Service) GetInfo() (*info.Info, error) {
	pages, err := s.pageService.Count()

	if err != nil {
		return nil, err
	}

	return &info.Info{
		PageCount: pages,
	}, nil
}
