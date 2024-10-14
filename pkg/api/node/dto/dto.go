package dto

import (
	"juno/pkg/api/node"

	"github.com/google/uuid"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Node struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func NewNodeFromDomain(n *node.Node) *Node {
	return &Node{
		ID:      n.ID.String(),
		OwnerID: n.OwnerID.String(),
		Address: n.Address,
		Shards:  n.Shards,
	}
}

func (n Node) ToDomain() (*node.Node, error) {
	id, err := uuid.Parse(n.ID)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.Parse(n.OwnerID)
	if err != nil {
		return nil, err
	}
	return node.New(id, ownerID, n.Address, n.Shards), nil
}

type GetNodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error
	Node    *Node  `json:"node,omitempty"`    // Only present when successful
}

func NewSuccessGetNodeResponse(node *node.Node) GetNodeResponse {
	n := NewNodeFromDomain(node)
	return GetNodeResponse{
		Status: SUCCESS,
		Node:   n,
	}
}

func NewErrorGetNodeResponse(message string) GetNodeResponse {
	return GetNodeResponse{
		Status:  ERROR,
		Message: message,
	}
}

type CreateNodeRequest struct {
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func (r CreateNodeRequest) ToDomain() node.Node {
	return node.Node{
		ID:      uuid.New(),
		Address: r.Address,
		Shards:  r.Shards,
	}
}

type CreateNodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Node *Node `json:"node,omitempty"` // Only present when successful
}

func NewSuccessCreateNodeResponse(node *node.Node) CreateNodeResponse {
	n := NewNodeFromDomain(node)
	return CreateNodeResponse{
		Status: SUCCESS,
		Node:   n,
	}
}

func NewErrorCreateNodeResponse(message string) CreateNodeResponse {
	return CreateNodeResponse{
		Status:  ERROR,
		Message: message,
	}
}

type UpdateNodeRequest struct {
	Address string `json:"address"`
	Shards  []int  `json:"shards"`
}

func (r UpdateNodeRequest) ToDomain() (*node.Node, error) {
	return &node.Node{
		Address: r.Address,
		Shards:  r.Shards,
	}, nil
}

type UpdateNodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Node *Node `json:"node,omitempty"` // Only present when successful
}

func NewSuccessUpdateNodeResponse(node *node.Node) UpdateNodeResponse {
	return UpdateNodeResponse{
		Status: SUCCESS,
		Node:   NewNodeFromDomain(node),
	}
}

func NewErrorUpdateNodeResponse(message string) UpdateNodeResponse {
	return UpdateNodeResponse{
		Status:  ERROR,
		Message: message,
	}
}

type DeleteNodeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error
}

func NewSuccessDeleteNodeResponse() DeleteNodeResponse {
	return DeleteNodeResponse{
		Status: SUCCESS,
	}
}

func NewErrorDeleteNodeResponse(message string) DeleteNodeResponse {
	return DeleteNodeResponse{
		Status:  ERROR,
		Message: message,
	}
}
