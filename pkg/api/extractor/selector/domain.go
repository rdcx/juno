package selector

import (
	"context"
	"errors"
	"juno/pkg/can"
	"juno/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("selector not found")

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type Selector struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string
	Value      string
	Visibility Visibility
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s Selector) Validate() error {
	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if s.Value == "" {
		errs = append(errs, errors.New("value is required"))
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	return nil
}

type Service interface {
	Create(userID uuid.UUID, name string, value string, vis Visibility) (*Selector, error)
	Get(id uuid.UUID) (*Selector, error)
	ListByUserID(userID uuid.UUID) ([]*Selector, error)
}

type Repository interface {
	Create(selector *Selector) error
	Get(id uuid.UUID) (*Selector, error)
	ListByUserID(userID uuid.UUID) ([]*Selector, error)
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanRead(ctx context.Context, selector *Selector) can.Result
	CanUpdate(ctx context.Context, selector *Selector) can.Result
	CanDelete(ctx context.Context, selector *Selector) can.Result
	CanList(ctx context.Context, selectors []*Selector) can.Result
}
