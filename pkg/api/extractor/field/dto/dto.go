package dto

import (
	"juno/pkg/api/extractor/field"
	"time"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Field struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	SelectorID string `json:"selector_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func NewFieldFromDomain(s *field.Field) *Field {
	return &Field{
		ID:         s.ID.String(),
		Name:       s.Name,
		SelectorID: s.SelectorID.String(),
		Type:       string(s.Type),
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  s.UpdatedAt.Format(time.RFC3339),
	}
}

type CreateFieldRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	SelectorID string `json:"selector_id" binding:"required" validate:"uuid"`
}

type CreateFieldResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Field *Field `json:"field,omitempty"`
}

func NewSuccessCreateFieldResponse(s *field.Field) *CreateFieldResponse {
	return &CreateFieldResponse{
		Status: SUCCESS,
		Field:  NewFieldFromDomain(s),
	}
}

func NewErrorCreateFieldResponse(err error) *CreateFieldResponse {
	return &CreateFieldResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type GetFieldResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Field *Field `json:"field,omitempty"`
}

func NewSuccessGetFieldResponse(s *field.Field) *GetFieldResponse {
	return &GetFieldResponse{
		Status: SUCCESS,
		Field:  NewFieldFromDomain(s),
	}
}

func NewErrorGetFieldResponse(err error) *GetFieldResponse {
	return &GetFieldResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type ListFieldResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Fields []*Field `json:"fields,omitempty"`
}

func NewSuccessListFieldResponse(fields []*field.Field) *ListFieldResponse {
	var sels []*Field

	for _, s := range fields {
		sels = append(sels, NewFieldFromDomain(s))
	}

	return &ListFieldResponse{
		Status: SUCCESS,
		Fields: sels,
	}
}

func NewErrorListFieldResponse(err error) *ListFieldResponse {
	return &ListFieldResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
