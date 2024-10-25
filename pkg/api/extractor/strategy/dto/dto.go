package dto

import (
	"juno/pkg/api/extractor/strategy"
	"time"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Strategy struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewStrategyFromDomain(s *strategy.Strategy) *Strategy {
	return &Strategy{
		ID:        s.ID.String(),
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}

type CreateStrategyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateStrategyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Strategy *Strategy `json:"strategy,omitempty"`
}

func NewSuccessCreateStrategyResponse(s *strategy.Strategy) *CreateStrategyResponse {
	return &CreateStrategyResponse{
		Status:   SUCCESS,
		Strategy: NewStrategyFromDomain(s),
	}
}

func NewErrorCreateStrategyResponse(err error) *CreateStrategyResponse {
	return &CreateStrategyResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type GetStrategyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Strategy *Strategy `json:"strategy,omitempty"`
}

func NewSuccessGetStrategyResponse(s *strategy.Strategy) *GetStrategyResponse {
	return &GetStrategyResponse{
		Status:   SUCCESS,
		Strategy: NewStrategyFromDomain(s),
	}
}

func NewErrorGetStrategyResponse(err error) *GetStrategyResponse {
	return &GetStrategyResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type ListStrategyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Strategys []*Strategy `json:"strategys,omitempty"`
}

func NewSuccessListStrategyResponse(strategys []*strategy.Strategy) *ListStrategyResponse {
	var sels []*Strategy

	for _, s := range strategys {
		sels = append(sels, NewStrategyFromDomain(s))
	}

	return &ListStrategyResponse{
		Status:    SUCCESS,
		Strategys: sels,
	}
}

func NewErrorListStrategyResponse(err error) *ListStrategyResponse {
	return &ListStrategyResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type AddSelectorRequest struct {
	SelectorID string `json:"selector_id" binding:"required" validate:"uuid"`
}

type AddSelectorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessAddSelectorResponse() *AddSelectorResponse {
	return &AddSelectorResponse{
		Status: SUCCESS,
	}
}

func NewErrorAddSelectorResponse(err error) *AddSelectorResponse {
	return &AddSelectorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type RemoveSelectorRequest struct {
	SelectorID string `json:"selector_id" binding:"required" validate:"uuid"`
}

type RemoveSelectorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessRemoveSelectorResponse() *RemoveSelectorResponse {
	return &RemoveSelectorResponse{
		Status: SUCCESS,
	}
}

func NewErrorRemoveSelectorResponse(err error) *RemoveSelectorResponse {
	return &RemoveSelectorResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type AddFilterRequest struct {
	FilterID string `json:"filter_id" binding:"required" validate:"uuid"`
}

type AddFilterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessAddFilterResponse() *AddFilterResponse {
	return &AddFilterResponse{
		Status: SUCCESS,
	}
}

func NewErrorAddFilterResponse(err error) *AddFilterResponse {
	return &AddFilterResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type RemoveFilterRequest struct {
	FilterID string `json:"filter_id" binding:"required" validate:"uuid"`
}

type RemoveFilterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessRemoveFilterResponse() *RemoveFilterResponse {
	return &RemoveFilterResponse{
		Status: SUCCESS,
	}
}

func NewErrorRemoveFilterResponse(err error) *RemoveFilterResponse {
	return &RemoveFilterResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type AddFieldRequest struct {
	FieldID string `json:"field_id" binding:"required" validate:"uuid"`
}

type AddFieldResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessAddFieldResponse() *AddFieldResponse {
	return &AddFieldResponse{
		Status: SUCCESS,
	}
}

func NewErrorAddFieldResponse(err error) *AddFieldResponse {
	return &AddFieldResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type RemoveFieldRequest struct {
	FieldID string `json:"field_id" binding:"required" validate:"uuid"`
}

type RemoveFieldResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessRemoveFieldResponse() *RemoveFieldResponse {
	return &RemoveFieldResponse{
		Status: SUCCESS,
	}
}

func NewErrorRemoveFieldResponse(err error) *RemoveFieldResponse {
	return &RemoveFieldResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
