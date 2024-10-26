package handler

import (
	"bytes"
	"encoding/json"
	extractionDto "juno/pkg/node/extraction/dto"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type mockService struct{}

func (m *mockService) Extract(req extractionDto.ExtractionRequest) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"page_title": "test",
		},
	}, nil
}

func TestExtract(t *testing.T) {
	h := New(logrus.New(), &mockService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := extractionDto.ExtractionRequest{
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
	}

	encoded, err := json.Marshal(req)

	if err != nil {
		t.Fatal(err)
	}

	c.Request, _ = http.NewRequest(http.MethodPost, "/extract", bytes.NewBuffer(encoded))

	h.Extract(c)

	if c.Writer.Status() != http.StatusOK {
		t.Fatalf("unexpected status code: %d", c.Writer.Status())
	}

	var res extractionDto.ExtractionResponse

	err = json.NewDecoder(w.Body).Decode(&res)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Extractions) != 1 {
		t.Fatalf("unexpected data length: %d", len(res.Extractions))
	}

	if res.Extractions[0]["page_title"] != "test" {
		t.Fatalf("unexpected data: %v", res.Extractions)
	}
}
