package assignment

import (
	"context"
	"errors"
	"juno/pkg/can"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("assignment not found")
var ErrInternal = errors.New("internal error")

type Assignment struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	EntityID uuid.UUID
	Offset   int
	Length   int
}

type Repository interface {
	Get(id uuid.UUID) (*Assignment, error)
	ListByEntityID(entityID uuid.UUID) ([]*Assignment, error)
	Create(assignment *Assignment) error
	Update(assignment *Assignment) error
	Delete(id uuid.UUID) error
}

type Service interface {
	Get(id uuid.UUID) (*Assignment, error)
	ListByEntityID(entityID uuid.UUID) ([]*Assignment, error)
	Create(ownerID, entityID uuid.UUID, offset, length int) (*Assignment, error)
	Update(id uuid.UUID, offset, length int) (*Assignment, error)
	Delete(id uuid.UUID) error
}

type Handler interface {
	ListByEntityID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanList(ctx context.Context, as []*Assignment) can.Result
	CanRead(ctx context.Context, a *Assignment) can.Result
	CanUpdate(ctx context.Context, a *Assignment) can.Result
	CanDelete(ctx context.Context, a *Assignment) can.Result
}
