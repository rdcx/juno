package ranag

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
var ErrNotFound = errors.New("ranag not found")
var ErrInternal = errors.New("internal error")

type Repository interface {
	All() ([]*Ranag, error)
	Create(n *Ranag) error
	Get(id uuid.UUID) (*Ranag, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Ranag, error)
	FirstWhereAddress(address string) (*Ranag, error)
	Update(n *Ranag) error
	Delete(id uuid.UUID) error
}

type Service interface {
	AllShardsRanags() (map[int][]*Ranag, error)
	GroupByRange() (map[[2]int][]*Ranag, error)
	Get(id uuid.UUID) (*Ranag, error)
	ListByOwnerID(ownerID uuid.UUID) ([]*Ranag, error)
	Create(ownerID uuid.UUID, addr string, shardAssignments [][2]int) (*Ranag, error)
	Update(id uuid.UUID, n *Ranag) (*Ranag, error)
	Delete(id uuid.UUID) error
}

type Handler interface {
	AllShardsRanags(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanUpdate(ctx context.Context, n *Ranag) can.Result
	CanRead(ctx context.Context, n *Ranag) can.Result
	CanList(ctx context.Context, ranags []*Ranag) can.Result
	CanDelete(ctx context.Context, n *Ranag) can.Result
}

type Ranag struct {
	ID               uuid.UUID `json:"id"`
	OwnerID          uuid.UUID `json:"owner_id"`
	Address          string    `json:"address"`
	ShardAssignments [][2]int  `json:"shard_assignments"`
}

func New(id, ownerID uuid.UUID, address string, shardAssignments [][2]int) *Ranag {
	return &Ranag{
		ID:               id,
		OwnerID:          ownerID,
		Address:          address,
		ShardAssignments: shardAssignments,
	}
}
