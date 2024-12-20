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
	ID               string   `json:"id"`
	OwnerID          string   `json:"owner_id"`
	Address          string   `json:"address"`
	Status           string   `json:"status"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func NewNodeFromDomain(n *node.Node) *Node {
	return &Node{
		ID:               n.ID.String(),
		OwnerID:          n.OwnerID.String(),
		Address:          n.Address,
		ShardAssignments: n.ShardAssignments,
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
	return node.New(id, ownerID, n.Address, n.ShardAssignments), nil
}

type ListNodesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Nodes []*Node `json:"nodes,omitempty"` // Only present when successful
}

func NewSuccessListNodesResponse(nodes []*node.Node) ListNodesResponse {
	var n []*Node
	for _, node := range nodes {
		n = append(n, NewNodeFromDomain(node))
	}
	return ListNodesResponse{
		Status: SUCCESS,
		Nodes:  n,
	}
}

func NewErrorListNodesResponse(message string) ListNodesResponse {
	return ListNodesResponse{
		Status:  ERROR,
		Message: message,
	}
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
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r CreateNodeRequest) ToDomain() node.Node {
	return node.Node{
		ID:               uuid.New(),
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
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
	Address          string   `json:"address"`
	ShardAssignments [][2]int `json:"shard_assignments"`
}

func (r UpdateNodeRequest) ToDomain() (*node.Node, error) {
	return &node.Node{
		Address:          r.Address,
		ShardAssignments: r.ShardAssignments,
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

type AllShardsNodesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Shards map[int][]string `json:"shards,omitempty"` // Only present when successful
}

func NewSuccessAllShardsNodesResponse(shards map[int][]*node.Node) AllShardsNodesResponse {
	m := make(map[int][]string)
	for i, nodes := range shards {
		for _, node := range nodes {
			m[i] = append(m[i], node.Address)
		}
	}
	return AllShardsNodesResponse{
		Status: SUCCESS,
		Shards: m,
	}
}

func NewErrorAllShardsNodesResponse(message string) AllShardsNodesResponse {
	return AllShardsNodesResponse{
		Status:  ERROR,
		Message: message,
	}
}

type QueryResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Only present when there's an error

	Extraction interface{} `json:"extraction,omitempty"` // Only present when successful
}
