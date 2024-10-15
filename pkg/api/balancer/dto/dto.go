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
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func NewBalancerFromDomain(n *balancer.Balancer) *Balancer {
	return &Balancer{
		ID:      n.ID.String(),
		OwnerID: n.OwnerID.String(),
		Address: n.Address,
		Shards:  n.Shards,
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
	return balancer.New(id, ownerID, n.Address, n.Shards), nil
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
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func (r CreateBalancerRequest) ToDomain() balancer.Balancer {
	return balancer.Balancer{
		ID:      uuid.New(),
		Address: r.Address,
		Shards:  r.Shards,
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
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func (r UpdateBalancerRequest) ToDomain() (*balancer.Balancer, error) {
	return &balancer.Balancer{
		Address: r.Address,
		Shards:  r.Shards,
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
