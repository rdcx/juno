package dto

import (
	"juno/pkg/api/node"

	"github.com/google/uuid"
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
		Status: "success",
		Node:   n,
	}
}

func NewErrorGetNodeResponse(message string) GetNodeResponse {
	return GetNodeResponse{
		Status:  "error",
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
		Status: "success",
		Node:   n,
	}
}

func NewErrorCreateNodeResponse(message string) CreateNodeResponse {
	return CreateNodeResponse{
		Status:  "error",
		Message: message,
	}
}
