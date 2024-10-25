package strategy

import (
	"context"
	"errors"
	"juno/pkg/api/extractor/field"
	"juno/pkg/api/extractor/filter"
	"juno/pkg/api/extractor/selector"
	"juno/pkg/can"
	"juno/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("strategy not found")

type Strategy struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Selectors []*selector.Selector
	Filters   []*filter.Filter
	Fields    []*field.Field
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s Strategy) Validate() error {
	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if s.UserID == uuid.Nil {
		errs = append(errs, errors.New("user_id is required"))
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	return nil
}

type Service interface {
	Create(userID, name string) (*Strategy, error)
	Get(id uuid.UUID) (*Strategy, error)
	AddSelector(id, selectorID uuid.UUID) error
	AddFilter(id, filterID uuid.UUID) error
	AddField(id, fieldID uuid.UUID) error
	ListByUserID(userID uuid.UUID) ([]*Strategy, error)
}

type Repository interface {
	Create(strategy *Strategy) error
	Get(id uuid.UUID) (*Strategy, error)
	ListByUserID(userID uuid.UUID) ([]*Strategy, error)
	ListBySelectorID(selectorID uuid.UUID) ([]*Strategy, error)
}

type StrategySelectorRepository interface {
	AddSelector(strategyID, selectorID uuid.UUID) error
	ListSelectorIDs(strategyID uuid.UUID) ([]*selector.Selector, error)
	RemoveSelector(strategyID, selectorID uuid.UUID) error
}

type StrategyFilterRepository interface {
	AddFilter(strategyID, filterID uuid.UUID) error
	ListFilterIDs(strategyID uuid.UUID) ([]*filter.Filter, error)
	RemoveFilter(strategyID, filterID uuid.UUID) error
}

type StrategyFieldRepository interface {
	AddField(strategyID, fieldID uuid.UUID) error
	ListFieldIDs(strategyID uuid.UUID) ([]*field.Field, error)
	RemoveField(strategyID, fieldID uuid.UUID) error
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	ListBySelectorID(c *gin.Context)
}

type Policy interface {
	CanCreate() can.Result
	CanRead(ctx context.Context, strategy *Strategy) can.Result
	CanUpdate(ctx context.Context, strategy *Strategy) can.Result
	CanDelete(ctx context.Context, strategy *Strategy) can.Result
	CanList(ctx context.Context, strategys []*Strategy) can.Result
}
