package node

import (
	"context"
	"errors"
	"juno/pkg/can"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrAddressExists = errors.New("address already exists")
var ErrInvalidAddress = errors.New("invalid address")
var ErrInvalidShards = errors.New("invalid shards")
var ErrNotFound = errors.New("node not found")
var ErrInternal = errors.New("internal error")

type Repository interface {
	All() ([]*Node, error)
	Create(n *Node) error
	Get(id uuid.UUID) (*Node, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Node, error)
	FirstWhereAddress(address string) (*Node, error)
	Update(n *Node) error
	Delete(id uuid.UUID) error
}

type Service interface {
	AllShardsNodes() (map[int][]*Node, error)
	Get(id uuid.UUID) (*Node, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Node, error)
	Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*Node, error)
	Update(id uuid.UUID, n *Node) (*Node, error)
	Delete(id uuid.UUID) error
}

type Handler interface {
	AllShardsNodes(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanUpdate(ctx context.Context, n *Node) can.Result
	CanRead(ctx context.Context, n *Node) can.Result
	CanList(ctx context.Context, nodes []*Node) can.Result
	CanDelete(ctx context.Context, n *Node) can.Result
}

type Node struct {
	ID               uuid.UUID `json:"id"`
	OwnerID          uuid.UUID `json:"owner_id"`
	Address          string    `json:"address"`
	ShardAssignments [][2]int  `json:"shard_assignments"`
}

func New(id, ownerID uuid.UUID, address string, shardAssignments [][2]int) *Node {
	return &Node{
		ID:               id,
		OwnerID:          ownerID,
		Address:          address,
		ShardAssignments: shardAssignments,
	}
}
