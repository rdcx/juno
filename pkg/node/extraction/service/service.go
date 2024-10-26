package service

import (
	"errors"
	"juno/pkg/node/html"
	"juno/pkg/node/page"
	"juno/pkg/node/storage"

	extractionDto "juno/pkg/node/extraction/dto"

	"github.com/sirupsen/logrus"
)

type Service struct {
	logger         *logrus.Logger
	pageService    page.Service
	storageService storage.Service
	htmlService    html.Service
}

func New(
	logger *logrus.Logger,
	pageService page.Service,
	storageService storage.Service,
	htmlService html.Service,
) *Service {
	return &Service{
		logger:         logger,
		pageService:    pageService,
		storageService: storageService,
		htmlService:    htmlService,
	}
}

func getSelector(selectorID string, selectors []*extractionDto.Selector) (*extractionDto.Selector, error) {
	for _, s := range selectors {
		if s.ID == selectorID {
			return s, nil
		}
	}

	return nil, errors.New("field not found")
}

func allFieldsEmpty(data map[string]interface{}) bool {
	for _, v := range data {
		if v != "" {
			return false
		}
	}

	return true
}

func (s *Service) Extract(req extractionDto.ExtractionRequest) ([]map[string]interface{}, error) {
	extractions := make([]map[string]interface{}, 0)
	s.pageService.Iterator(func(p *page.Page) {
		for _, v := range p.Versions {
			body, err := s.storageService.Read(v.Hash)

			if err != nil {
				s.logger.WithError(err).Error("failed to get data from storage")
				return
			}

			pageData := map[string]interface{}{}

			for _, e := range req.Fields {

				// get the corresponding selector
				selector, err := getSelector(e.SelectorID, req.Selectors)

				if err != nil {
					s.logger.WithError(err).Error("failed to get selector")
					return
				}

				// get the value from the data
				val, err := s.htmlService.GetSelectorValue(body, selector.Value)

				if err != nil {
					s.logger.WithError(err).Error("failed to get selector value")
					return
				}

				pageData[e.Name] = val
			}

			if err != nil {
				s.logger.WithError(err).Error("failed to get title from HTML")
				return
			}

			if allFieldsEmpty(pageData) {
				return
			}

			pageData["_juno_meta_url"] = p.URL

			extractions = append(extractions, pageData)
		}
	})

	return extractions, nil
}
