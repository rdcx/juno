package balancer

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
var ErrNotFound = errors.New("balancer not found")
var ErrInternal = errors.New("internal error")

type Repository interface {
	All() ([]*Balancer, error)
	Create(n *Balancer) error
	Get(id uuid.UUID) (*Balancer, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Balancer, error)
	FirstWhereAddress(address string) (*Balancer, error)
	Update(n *Balancer) error
	Delete(id uuid.UUID) error
}

type Service interface {
	AllShardsBalancers() (map[int][]*Balancer, error)
	Get(id uuid.UUID) (*Balancer, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Balancer, error)
	Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*Balancer, error)
	Update(id uuid.UUID, n *Balancer) (*Balancer, error)
	Delete(id uuid.UUID) error
}

type Handler interface {
	AllShardsBalancers(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanUpdate(ctx context.Context, n *Balancer) can.Result
	CanRead(ctx context.Context, n *Balancer) can.Result
	CanList(ctx context.Context, balancers []*Balancer) can.Result
	CanDelete(ctx context.Context, n *Balancer) can.Result
}

type Balancer struct {
	ID               uuid.UUID `json:"id"`
	OwnerID          uuid.UUID `json:"owner_id"`
	Address          string    `json:"address"`
	ShardAssignments [][2]int  `json:"shard_assignments"`
}

func New(id, ownerID uuid.UUID, address string, shardAssignments [][2]int) *Balancer {
	return &Balancer{
		ID:               id,
		OwnerID:          ownerID,
		Address:          address,
		ShardAssignments: shardAssignments,
	}
}
