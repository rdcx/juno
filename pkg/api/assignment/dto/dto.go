package dto

import "juno/pkg/api/assignment"

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Assignment struct {
	ID       string `json:"id"`
	OwnerID  string `json:"owner_id"`
	EntityID string `json:"entity_id"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
}

type CreateAssignmentRequest struct {
	OwnerID  string `json:"owner_id"`
	EntityID string `json:"entity_id"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
}

type CreateAssignmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Assignment *Assignment `json:"assignment,omitempty"`
}

func NewAssignmentFromDomain(a *assignment.Assignment) *Assignment {
	return &Assignment{
		ID:       a.ID.String(),
		OwnerID:  a.OwnerID.String(),
		EntityID: a.EntityID.String(),
		Offset:   a.Offset,
		Length:   a.Length,
	}
}

func NewSuccessCreateAssignmentResponse(a *assignment.Assignment) *CreateAssignmentResponse {
	return &CreateAssignmentResponse{
		Status:     SUCCESS,
		Assignment: NewAssignmentFromDomain(a),
	}
}

func NewErrorCreateAssignmentResponse(err error) *CreateAssignmentResponse {
	return &CreateAssignmentResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type GetAssignmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Assignment *Assignment `json:"assignment,omitempty"`
}

func NewSuccessGetAssignmentResponse(a *assignment.Assignment) *GetAssignmentResponse {
	return &GetAssignmentResponse{
		Status:     SUCCESS,
		Assignment: NewAssignmentFromDomain(a),
	}
}

func NewErrorGetAssignmentResponse(err error) *GetAssignmentResponse {
	return &GetAssignmentResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type ListAssignmentsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Assignments []*Assignment `json:"assignments,omitempty"`
}

func NewSuccessListAssignmentsResponse(assignments []*assignment.Assignment) *ListAssignmentsResponse {
	var a []*Assignment
	for _, assignment := range assignments {
		a = append(a, NewAssignmentFromDomain(assignment))
	}
	return &ListAssignmentsResponse{
		Status:      SUCCESS,
		Assignments: a,
	}
}

func NewErrorListAssignmentsResponse(err error) *ListAssignmentsResponse {
	return &ListAssignmentsResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type UpdateAssignmentRequest struct {
	EntityID string `json:"entity_id"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
}

type UpdateAssignmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Assignment *Assignment `json:"assignment,omitempty"`
}

func NewSuccessUpdateAssignmentResponse(a *assignment.Assignment) *UpdateAssignmentResponse {
	return &UpdateAssignmentResponse{
		Status:     SUCCESS,
		Assignment: NewAssignmentFromDomain(a),
	}
}

func NewErrorUpdateAssignmentResponse(err error) *UpdateAssignmentResponse {
	return &UpdateAssignmentResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}

type DeleteAssignmentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Assignment *Assignment `json:"assignment,omitempty"`
}

func NewSuccessDeleteAssignmentResponse(a *assignment.Assignment) *DeleteAssignmentResponse {
	return &DeleteAssignmentResponse{
		Status:     SUCCESS,
		Assignment: NewAssignmentFromDomain(a),
	}
}

func NewErrorDeleteAssignmentResponse(err error) *DeleteAssignmentResponse {
	return &DeleteAssignmentResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
