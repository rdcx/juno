package dto

import (
	fieldDto "juno/pkg/api/extractor/field/dto"
	filterDto "juno/pkg/api/extractor/filter/dto"
	selectorDto "juno/pkg/api/extractor/selector/dto"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type RangeAggregatorRequest struct {
	Selectors []*selectorDto.Selector `json:"selectors" binding:"required"`
	Fields    []*fieldDto.Field       `json:"fields" binding:"required"`
	Filters   []*filterDto.Filter     `json:"filters" binding:"required"`
}

type RangeAggregatorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Aggregations []map[string]interface{} `json:"aggregations,omitempty"`
}

func NewSuccessRangeAggregatorResponse(aggregations []map[string]interface{}) *RangeAggregatorResponse {
	return &RangeAggregatorResponse{
		Status:       SUCCESS,
		Aggregations: aggregations,
	}
}

func NewErrorRangeAggregatorResponse(err error) *RangeAggregatorResponse {
	return &RangeAggregatorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
