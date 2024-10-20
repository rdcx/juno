package service

import (
	monkeyService "juno/pkg/monkey/service"
	htmlService "juno/pkg/node/html/service"
	"juno/pkg/node/page"
	pageRepo "juno/pkg/node/page/repo/mem"
	pageService "juno/pkg/node/page/service"
	storageService "juno/pkg/node/storage/service"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTitles(t *testing.T) {

	logger := logrus.New()
	pageRepo := pageRepo.New()
	pageService := pageService.New(pageRepo)
	storageService := storageService.New(t.TempDir())
	htmlService := htmlService.New()
	monkeyService := monkeyService.New()

	s := New(
		logger,
		pageService,
		storageService,
		htmlService,
		monkeyService,
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

	titles, err := s.Titles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(titles) != 1 {
		t.Errorf("expected 1 title, got %d", len(titles))
	}

	if titles["http://example.com"] != "Test" {
		t.Errorf("expected title to be 'Test', got %s", titles["http://example.com"])
	}
}
