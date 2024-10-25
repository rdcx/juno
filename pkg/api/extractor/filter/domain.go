package filter

import (
	"context"
	"errors"
	"juno/pkg/can"
	"juno/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("filter not found")

type FilterType string

const (
	FilterTypeStringEquals   FilterType = "string_equals"
	FilterTypeStringContains FilterType = "string_contains"
)

type Filter struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      FilterType
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s Filter) Validate() error {
	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if s.Type == "" {
		errs = append(errs, errors.New("type is required"))
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
	Create(userID uuid.UUID, name string, t FilterType, value string) (*Filter, error)
	Get(id uuid.UUID) (*Filter, error)
	ListByUserID(userID uuid.UUID) ([]*Filter, error)
}

type Repository interface {
	Create(filter *Filter) error
	Get(id uuid.UUID) (*Filter, error)
	ListByUserID(userID uuid.UUID) ([]*Filter, error)
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanRead(ctx context.Context, filter *Filter) can.Result
	CanUpdate(ctx context.Context, filter *Filter) can.Result
	CanDelete(ctx context.Context, filter *Filter) can.Result
	CanList(ctx context.Context, filters []*Filter) can.Result
}
