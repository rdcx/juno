package dto

import (
	"juno/pkg/api/ranag"

	"github.com/google/uuid"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Ranag struct {
	ID               string   `json:"id"`
	OwnerID          string   `json:"owner_id"`
	Address          string   `json:"address"`
	Status           string   `json:"status"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func NewRanagFromDomain(n *ranag.Ranag) *Ranag {
	return &Ranag{
		ID:               n.ID.String(),
		OwnerID:          n.OwnerID.String(),
		Address:          n.Address,
		ShardAssignments: n.ShardAssignments,
	}
}

func (n Ranag) ToDomain() (*ranag.Ranag, error) {
	id, err := uuid.Parse(n.ID)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.Parse(n.OwnerID)
	if err != nil {
		return nil, err
	}
	return ranag.New(id, ownerID, n.Address, n.ShardAssignments), nil
}

type ListRanagsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Ranags []*Ranag `json:"ranags,omitempty"` // Only present when successful
}

func NewSuccessListRanagsResponse(ranags []*ranag.Ranag) ListRanagsResponse {
	var n []*Ranag
	for _, ranag := range ranags {
		n = append(n, NewRanagFromDomain(ranag))
	}
	return ListRanagsResponse{
		Status: SUCCESS,
		Ranags: n,
	}
}

func NewErrorListRanagsResponse(message string) ListRanagsResponse {
	return ListRanagsResponse{
		Status:  ERROR,
		Message: message,
	}
}

type GetRanagResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error
	Ranag   *Ranag `json:"ranag,omitempty"`   // Only present when successful
}

func NewSuccessGetRanagResponse(ranag *ranag.Ranag) GetRanagResponse {
	n := NewRanagFromDomain(ranag)
	return GetRanagResponse{
		Status: SUCCESS,
		Ranag:  n,
	}
}

func NewErrorGetRanagResponse(message string) GetRanagResponse {
	return GetRanagResponse{
		Status:  ERROR,
		Message: message,
	}
}

type CreateRanagRequest struct {
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r CreateRanagRequest) ToDomain() ranag.Ranag {
	return ranag.Ranag{
		ID:               uuid.New(),
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
	}
}

type CreateRanagResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Ranag *Ranag `json:"ranag,omitempty"` // Only present when successful
}

func NewSuccessCreateRanagResponse(ranag *ranag.Ranag) CreateRanagResponse {
	n := NewRanagFromDomain(ranag)
	return CreateRanagResponse{
		Status: SUCCESS,
		Ranag:  n,
	}
}

func NewErrorCreateRanagResponse(message string) CreateRanagResponse {
	return CreateRanagResponse{
		Status:  ERROR,
		Message: message,
	}
}

type UpdateRanagRequest struct {
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r UpdateRanagRequest) ToDomain() (*ranag.Ranag, error) {
	return &ranag.Ranag{
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
	}, nil
}

type UpdateRanagResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Ranag *Ranag `json:"ranag,omitempty"` // Only present when successful
}

func NewSuccessUpdateRanagResponse(ranag *ranag.Ranag) UpdateRanagResponse {
	return UpdateRanagResponse{
		Status: SUCCESS,
		Ranag:  NewRanagFromDomain(ranag),
	}
}

func NewErrorUpdateRanagResponse(message string) UpdateRanagResponse {
	return UpdateRanagResponse{
		Status:  ERROR,
		Message: message,
	}
}

type DeleteRanagResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error
}

func NewSuccessDeleteRanagResponse() DeleteRanagResponse {
	return DeleteRanagResponse{
		Status: SUCCESS,
	}
}

func NewErrorDeleteRanagResponse(message string) DeleteRanagResponse {
	return DeleteRanagResponse{
		Status:  ERROR,
		Message: message,
	}
}

type AllShardsRanagsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Shards map[int][]string `json:"shards,omitempty"` // Only present when successful
}

func NewSuccessAllShardsRanagsResponse(shards map[int][]*ranag.Ranag) AllShardsRanagsResponse {
	m := make(map[int][]string)
	for i, ranags := range shards {
		for _, ranag := range ranags {
			m[i] = append(m[i], ranag.Address)
		}
	}
	return AllShardsRanagsResponse{
		Status: SUCCESS,
		Shards: m,
	}
}

func NewErrorAllShardsRanagsResponse(message string) AllShardsRanagsResponse {
	return AllShardsRanagsResponse{
		Status:  ERROR,
		Message: message,
	}
}
