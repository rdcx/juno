package dto

import (
	"juno/pkg/api/extractor/filter"
	"time"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Filter struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewFilterFromDomain(s *filter.Filter) *Filter {
	return &Filter{
		ID:        s.ID.String(),
		Name:      s.Name,
		Value:     s.Value,
		Type:      string(s.Type),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}

type CreateFilterRequest struct {
	FieldID string `json:"field_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Value   string `json:"value" binding:"required"`
}

type CreateFilterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Filter *Filter `json:"filter,omitempty"`
}

func NewSuccessCreateFilterResponse(s *filter.Filter) *CreateFilterResponse {
	return &CreateFilterResponse{
		Status: SUCCESS,
		Filter: NewFilterFromDomain(s),
	}
}

func NewErrorCreateFilterResponse(err error) *CreateFilterResponse {
	return &CreateFilterResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type GetFilterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Filter *Filter `json:"filter,omitempty"`
}

func NewSuccessGetFilterResponse(s *filter.Filter) *GetFilterResponse {
	return &GetFilterResponse{
		Status: SUCCESS,
		Filter: NewFilterFromDomain(s),
	}
}

func NewErrorGetFilterResponse(err error) *GetFilterResponse {
	return &GetFilterResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type ListFilterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Filters []*Filter `json:"filters,omitempty"`
}

func NewSuccessListFilterResponse(filters []*filter.Filter) *ListFilterResponse {
	var sels []*Filter

	for _, s := range filters {
		sels = append(sels, NewFilterFromDomain(s))
	}

	return &ListFilterResponse{
		Status:  SUCCESS,
		Filters: sels,
	}
}

func NewErrorListFilterResponse(err error) *ListFilterResponse {
	return &ListFilterResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
