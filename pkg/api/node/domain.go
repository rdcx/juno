package node

import (
	"errors"
	"juno/pkg/api/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrAddressExists = errors.New("address already exists")
var ErrInvalidAddress = errors.New("invalid address")
var ErrInvalidShards = errors.New("invalid shards")
var ErrNotFound = errors.New("node not found")

type Repository interface {
	Create(n *Node) error
	Get(id uuid.UUID) (*Node, error)
	FirstWhereAddress(address string) (*Node, error)
	Update(n *Node) error
	Delete(id uuid.UUID) error
}

type Service interface {
	Get(u *user.User, id uuid.UUID) (*Node, error)
	Create(u *user.User, addr string, shards []int) (*Node, error)
	Update(u *user.User, n *Node) error
	Delete(u *user.User, id uuid.UUID) error
}

type Handler interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Node struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"owner_id"`
	Address string    `json:"address"`
	Shards  []int     `json:"shards"`
}

func New(id, ownerID uuid.UUID, address string, shards []int) *Node {
	return &Node{
		ID:      id,
		OwnerID: ownerID,
		Address: address,
		Shards:  shards,
	}
}
