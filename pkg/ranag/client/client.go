package client

import (
	"bytes"
	"encoding/json"
	"errors"
	fieldDto "juno/pkg/api/extractor/field/dto"
	filterDto "juno/pkg/api/extractor/filter/dto"
	selectorDto "juno/pkg/api/extractor/selector/dto"

	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/selector"
	"net/http"

	dto "juno/pkg/ranag/dto"
)

type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c Client) SendRangeAggregationRequest(selectors []*selector.Selector, fields []*field.Field, filters []*filter.Filter) (*dto.RangeAggregatorResponse, error) {

	selectorDtos := make([]*selectorDto.Selector, 0, len(selectors))

	for _, s := range selectors {
		selectorDtos = append(selectorDtos, selectorDto.NewSelectorFromDomain(s))
	}

	fieldDtos := make([]*fieldDto.Field, 0, len(fields))

	for _, f := range fields {
		fieldDtos = append(fieldDtos, fieldDto.NewFieldFromDomain(f))
	}

	filterDtos := make([]*filterDto.Filter, 0, len(filters))

	for _, f := range filters {
		filterDtos = append(filterDtos, filterDto.NewFilterFromDomain(f))
	}

	req := &dto.RangeAggregatorRequest{
		Selectors: selectorDtos,
		Fields:    fieldDtos,
		Filters:   filterDtos,
	}

	encoded, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://"+c.baseURL+"/aggregation", "application/json", bytes.NewBuffer(encoded))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	var rangeAggregatorResponse dto.RangeAggregatorResponse

	err = json.NewDecoder(resp.Body).Decode(&rangeAggregatorResponse)

	if err != nil {
		return nil, err
	}

	return &rangeAggregatorResponse, nil
}
