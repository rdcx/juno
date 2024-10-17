package service

import (
	"juno/pkg/monkey"
	"juno/pkg/node/html"
	"juno/pkg/node/page"
	"juno/pkg/node/storage"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger         *logrus.Logger
	pageService    page.Service
	storageService storage.Service
	htmlService    html.Service
	monkeyService  monkey.Service
}

func New(
	logger *logrus.Logger,
	pageService page.Service,
	storageService storage.Service,
	htmlService html.Service,
	monkeyService monkey.Service,
) *Service {
	return &Service{
		logger:         logger,
		pageService:    pageService,
		storageService: storageService,
		htmlService:    htmlService,
		monkeyService:  monkeyService,
	}
}

func (s *Service) Titles() ([]string, error) {
	var titles []string
	s.pageService.Iterator(func(p *page.Page) {
		for _, v := range p.Versions {
			data, err := s.storageService.Read(v.Hash)

			if err != nil {
				s.logger.WithError(err).Error("failed to get data from storage")
				return
			}

			title, err := s.htmlService.Title(data)

			if err != nil {
				s.logger.WithError(err).Error("failed to get title from HTML")
				return
			}

			titles = append(titles, title)
		}
	})

	return titles, nil
}

func (s *Service) Execute(src string) ([]byte, error) {
	return s.monkeyService.Execute(src)
}
