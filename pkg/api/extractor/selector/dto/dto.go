package dto

import (
	"juno/pkg/api/extractor/selector"
	"time"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Selector struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func NewSelectorFromDomain(s *selector.Selector) *Selector {
	return &Selector{
		ID:         s.ID.String(),
		Name:       s.Name,
		Value:      s.Value,
		Visibility: string(s.Visibility),
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  s.UpdatedAt.Format(time.RFC3339),
	}
}

type CreateSelectorRequest struct {
	Name       string `json:"name" binding:"required"`
	Value      string `json:"value" binding:"required"`
	Visibility string `json:"visibility" binding:"required"`
}

type CreateSelectorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Selector *Selector `json:"selector,omitempty"`
}

func NewSuccessCreateSelectorResponse(s *selector.Selector) *CreateSelectorResponse {
	return &CreateSelectorResponse{
		Status:   SUCCESS,
		Selector: NewSelectorFromDomain(s),
	}
}

func NewErrorCreateSelectorResponse(err error) *CreateSelectorResponse {
	return &CreateSelectorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type GetSelectorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Selector *Selector `json:"selector,omitempty"`
}

func NewSuccessGetSelectorResponse(s *selector.Selector) *GetSelectorResponse {
	return &GetSelectorResponse{
		Status:   SUCCESS,
		Selector: NewSelectorFromDomain(s),
	}
}

func NewErrorGetSelectorResponse(err error) *GetSelectorResponse {
	return &GetSelectorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type ListSelectorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Selectors []*Selector `json:"selectors,omitempty"`
}

func NewSuccessListSelectorResponse(selectors []*selector.Selector) *ListSelectorResponse {
	var sels []*Selector

	for _, s := range selectors {
		sels = append(sels, NewSelectorFromDomain(s))
	}

	return &ListSelectorResponse{
		Status:    SUCCESS,
		Selectors: sels,
	}
}

func NewErrorListSelectorResponse(err error) *ListSelectorResponse {
	return &ListSelectorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
