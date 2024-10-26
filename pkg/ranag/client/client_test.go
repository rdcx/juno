package client

import (
	"testing"

	"github.com/google/uuid"
	"github.com/h2non/gock"

	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/ranag/dto"
)

func TestSendRangeAggregationRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("/aggregation").
			Reply(200).
			JSON(dto.NewSuccessRangeAggregatorResponse([]map[string]interface{}{
				{
					"product_title": "charger",
					"price":         100,
				},
			}))

		client := New("localhost:8080")
		selectors := []*selector.Selector{
			{
				Value: "#productTitle",
			},
		}

		fields := []*field.Field{
			{
				SelectorID: uuid.New(),
			},
		}

		filters := []*filter.Filter{
			{
				Name:    "test",
				FieldID: uuid.New(),
				Type:    "string_equals",
				Value:   "test",
			},
		}

		resp, err := client.SendRangeAggregationRequest(
			selectors,
			fields,
			filters,
		)

		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})

	t.Run("error", func(t *testing.T) {

		defer gock.Off()

		gock.New("http://localhost:8080").
			Post("/aggregation").
			Reply(500)

		client := New("localhost:8080")

		_, err := client.SendRangeAggregationRequest(nil, nil, nil)

		if err == nil {
			t.Fatal("expected error")
		}
	})
}
