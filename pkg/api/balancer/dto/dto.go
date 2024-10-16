package dto

import (
	"juno/pkg/api/balancer"

	"github.com/google/uuid"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Balancer struct {
	ID               string   `json:"id"`
	OwnerID          string   `json:"owner_id"`
	Address          string   `json:"address"`
	Status           string   `json:"status"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func NewBalancerFromDomain(n *balancer.Balancer) *Balancer {
	return &Balancer{
		ID:               n.ID.String(),
		OwnerID:          n.OwnerID.String(),
		Address:          n.Address,
		ShardAssignments: n.ShardAssignments,
	}
}

func (n Balancer) ToDomain() (*balancer.Balancer, error) {
	id, err := uuid.Parse(n.ID)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.Parse(n.OwnerID)
	if err != nil {
		return nil, err
	}
	return balancer.New(id, ownerID, n.Address, n.ShardAssignments), nil
}

type ListBalancersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Balancers []*Balancer `json:"balancers,omitempty"` // Only present when successful
}

func NewSuccessListBalancersResponse(balancers []*balancer.Balancer) ListBalancersResponse {
	var n []*Balancer
	for _, balancer := range balancers {
		n = append(n, NewBalancerFromDomain(balancer))
	}
	return ListBalancersResponse{
		Status:    SUCCESS,
		Balancers: n,
	}
}

func NewErrorListBalancersResponse(message string) ListBalancersResponse {
	return ListBalancersResponse{
		Status:  ERROR,
		Message: message,
	}
}

type GetBalancerResponse struct {
	Status   string    `json:"status"`
	Message  string    `json:"message,omitempty"`  // Only present when there's an error
	Balancer *Balancer `json:"balancer,omitempty"` // Only present when successful
}

func NewSuccessGetBalancerResponse(balancer *balancer.Balancer) GetBalancerResponse {
	n := NewBalancerFromDomain(balancer)
	return GetBalancerResponse{
		Status:   SUCCESS,
		Balancer: n,
	}
}

func NewErrorGetBalancerResponse(message string) GetBalancerResponse {
	return GetBalancerResponse{
		Status:  ERROR,
		Message: message,
	}
}

type CreateBalancerRequest struct {
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r CreateBalancerRequest) ToDomain() balancer.Balancer {
	return balancer.Balancer{
		ID:               uuid.New(),
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
	}
}

type CreateBalancerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Balancer *Balancer `json:"balancer,omitempty"` // Only present when successful
}

func NewSuccessCreateBalancerResponse(balancer *balancer.Balancer) CreateBalancerResponse {
	n := NewBalancerFromDomain(balancer)
	return CreateBalancerResponse{
		Status:   SUCCESS,
		Balancer: n,
	}
}

func NewErrorCreateBalancerResponse(message string) CreateBalancerResponse {
	return CreateBalancerResponse{
		Status:  ERROR,
		Message: message,
	}
}

type UpdateBalancerRequest struct {
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r UpdateBalancerRequest) ToDomain() (*balancer.Balancer, error) {
	return &balancer.Balancer{
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
	}, nil
}

type UpdateBalancerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Balancer *Balancer `json:"balancer,omitempty"` // Only present when successful
}

func NewSuccessUpdateBalancerResponse(balancer *balancer.Balancer) UpdateBalancerResponse {
	return UpdateBalancerResponse{
		Status:   SUCCESS,
		Balancer: NewBalancerFromDomain(balancer),
	}
}

func NewErrorUpdateBalancerResponse(message string) UpdateBalancerResponse {
	return UpdateBalancerResponse{
		Status:  ERROR,
		Message: message,
	}
}

type DeleteBalancerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error
}

func NewSuccessDeleteBalancerResponse() DeleteBalancerResponse {
	return DeleteBalancerResponse{
		Status: SUCCESS,
	}
}

func NewErrorDeleteBalancerResponse(message string) DeleteBalancerResponse {
	return DeleteBalancerResponse{
		Status:  ERROR,
		Message: message,
	}
}

type AllShardsBalancersResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Shards map[int][]string `json:"shards,omitempty"` // Only present when successful
}

func NewSuccessAllShardsBalancersResponse(shards map[int][]*balancer.Balancer) AllShardsBalancersResponse {
	m := make(map[int][]string)
	for i, balancers := range shards {
		for _, balancer := range balancers {
			m[i] = append(m[i], balancer.Address)
		}
	}
	return AllShardsBalancersResponse{
		Status: SUCCESS,
		Shards: m,
	}
}

func NewErrorAllShardsBalancersResponse(message string) AllShardsBalancersResponse {
	return AllShardsBalancersResponse{
		Status:  ERROR,
		Message: message,
	}
}
