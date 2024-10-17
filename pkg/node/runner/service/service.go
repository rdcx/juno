package service

import (
	"juno/pkg/monkey"
	"juno/pkg/node/page"
	"juno/pkg/node/storage"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger         *logrus.Logger
	pageService    page.Service
	storageService storage.Service
	monkeyService  monkey.Service
}

func New(
	logger *logrus.Logger,
	pageService page.Service,
	storageService storage.Service,
	monkeyService monkey.Service,
) *Service {
	return &Service{
		logger:         logger,
		pageService:    pageService,
		storageService: storageService,
		monkeyService:  monkeyService,
	}
}

func (s *Service) Execute(src string) ([]byte, error) {
	return s.monkeyService.Execute(src)
}
