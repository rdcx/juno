package field

import (
	"context"
	"errors"
	"juno/pkg/can"
	"juno/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("field not found")

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeFloat   FieldType = "float"
)

type Field struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	SelectorID uuid.UUID
	Name       string
	Type       FieldType
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s Field) Validate() error {
	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if s.UserID == uuid.Nil {
		errs = append(errs, errors.New("user_id is required"))
	}

	if s.SelectorID == uuid.Nil {
		errs = append(errs, errors.New("selector_id is required"))
	}

	if s.Type == "" {
		errs = append(errs, errors.New("type is required"))
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	return nil
}

type Service interface {
	Create(userID, selectorID uuid.UUID, name string, t FieldType) (*Field, error)
	Get(id uuid.UUID) (*Field, error)
	ListByUserID(userID uuid.UUID) ([]*Field, error)
	ListBySelectorID(selectorID uuid.UUID) ([]*Field, error)
}

type Repository interface {
	Create(field *Field) error
	Get(id uuid.UUID) (*Field, error)
	ListByUserID(userID uuid.UUID) ([]*Field, error)
	ListBySelectorID(selectorID uuid.UUID) ([]*Field, error)
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	ListBySelectorID(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanRead(ctx context.Context, field *Field) can.Result
	CanUpdate(ctx context.Context, field *Field) can.Result
	CanDelete(ctx context.Context, field *Field) can.Result
	CanList(ctx context.Context, fields []*Field) can.Result
}
