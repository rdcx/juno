package service

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) ExtractLinks(body []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})

	return links, nil
}

func (s *Service) Title(body []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))

	if err != nil {
		return "", err
	}

	return doc.Find("title").Text(), nil
}
