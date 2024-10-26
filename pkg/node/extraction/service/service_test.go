package service

import (
	htmlService "juno/pkg/node/html/service"
	"juno/pkg/node/page"
	pageRepo "juno/pkg/node/page/repo/mem"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"

	extractionDto "juno/pkg/node/extraction/dto"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestExtract(t *testing.T) {

	logger := logrus.New()
	pageRepo := pageRepo.New()
	pageService := pageService.New(pageRepo)
	storageService := storageService.New(t.TempDir())
	htmlService := htmlService.New()

	s := New(
		logger,
		pageService,
		storageService,
		htmlService,
	)

	body := []byte("<html><head><title>Test</title></head><body></body></html>")

	p := page.NewPage("http://example.com")
	err := pageService.Create(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vHash := page.NewVersionHash(body)

	pageService.AddVersion(p.ID, page.NewVersion(
		vHash,
	))

	err = storageService.Write(vHash, body)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := s.Extract(
		extractionDto.ExtractionRequest{
			Shard: 68735,
			Selectors: []*extractionDto.Selector{
				{
					ID:    "1",
					Value: "title",
				},
			},
			Fields: []*extractionDto.Field{
				{
					SelectorID: "1",
					Name:       "page_title",
				},
			},
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1, got %d", len(data))
	}

	if data[0]["page_title"] != "Test" {
		t.Fatalf("expected Test, got %s", data[0]["page_title"])
	}

	if data[0]["_juno_meta_url"] != "http://example.com" {
		t.Fatalf("expected http://example.com, got %s", data[0]["_juno_meta_url"])
	}
}
