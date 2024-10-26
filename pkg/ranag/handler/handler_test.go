package handler

import (
	"bytes"
	"encoding/json"
	"juno/pkg/ranag/dto"
	"net/http"
	"net/http/httptest"
	"testing"

	fieldDto "juno/pkg/api/extractor/field/dto"
	filterDto "juno/pkg/api/extractor/filter/dto"
	selectorDto "juno/pkg/api/extractor/selector/dto"

	"github.com/gin-gonic/gin"
)

type mockService struct{}

func (m *mockService) RangeAggregate(offset, total int, req dto.RangeAggregatorRequest) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"product_title": "test",
		},
	}, nil
}

func TestRangeAggregate(t *testing.T) {
	h := New(&mockService{})

	req := dto.RangeAggregatorRequest{
		Offset: 0,
		Total:  1,
		Selectors: []*selectorDto.Selector{
			{
				ID:   "1",
				Name: "#productTitle",
			},
		},
		Fields: []*fieldDto.Field{
			{
				ID:         "1",
				SelectorID: "1",
				Name:       "product_title",
			},
		},
		Filters: []*filterDto.Filter{
			{
				FieldID: "1",
				Type:    "string_equals",
				Value:   "test",
			},
		},
	}

	encoded, err := json.Marshal(req)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request, err = http.NewRequest(http.MethodPost, "/aggregate", bytes.NewReader(encoded))

	h.RangeAggregate(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var resp dto.RangeAggregatorResponse

	err = json.NewDecoder(w.Body).Decode(&resp)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != dto.SUCCESS {
		t.Errorf("Expected status %s, got %s", dto.SUCCESS, resp.Status)
	}

	if len(resp.Aggregations) != 1 {
		t.Fatalf("Expected 1 aggregation, got %d", len(resp.Aggregations))
	}

	aggregation := resp.Aggregations[0]

	if aggregation["product_title"] != "test" {
		t.Errorf("Expected product_title to be test, got %s", aggregation["product_title"])
	}
}
