package service

import "testing"

func TestExtractLinks(t *testing.T) {
	t.Run("should return empty slice when no links found", func(t *testing.T) {
		body := `<html><head><title>Test</title></head><body></body></html>`

		links, err := New().ExtractLinks([]byte(body))

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if len(links) != 0 {
			t.Errorf("expected 0 links but got %d", len(links))
		}
	})

	t.Run("should return links when found", func(t *testing.T) {
		body := `<html><head><title>Test</title></head><body><a href="http://example.com">Example</a> <a href="http://example.net">Net Example</a></body></html>`

		links, err := New().ExtractLinks([]byte(body))

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if len(links) != 2 {
			t.Errorf("expected 1 link but got %d", len(links))
		}

		if links[0] != "http://example.com" {
			t.Errorf("expected http://example.com but got %s", links[0])
		}

		if links[1] != "http://example.net" {
			t.Errorf("expected http://example.net but got %s", links[1])
		}
	})
}
