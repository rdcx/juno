package strategy

import (
	"errors"
	"juno/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StringMatchType string
type JobStatus string
type FilterType string

type StringMatch struct {
	Value     string          `json:"value"`
	MatchType StringMatchType `json:"type"`
}

type Filter struct {
	FilterType FilterType `json:"type"`
	Value      string     `json:"value"`
}

type Strategy struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Instruction
	Filters   map[string]*Filter `json:"filters"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func validateStrategy(name string, instruction Instruction, filters map[string]*Filter) error {

	var errs []error

	if name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if len(instruction.Selectors) == 0 {
		errs = append(errs, errors.New("selectors are required"))
	}

	if len(instruction.OutputFormat) == 0 {
		errs = append(errs, errors.New("output format is required"))
	}

	for _, filter := range filters {
		if filter.FilterType == "" {
			errs = append(errs, errors.New("filter type is required"))
		}
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	return nil
}

func NewStrategy(userID uuid.UUID, name string, instruction Instruction, filters map[string]*Filter) (*Strategy, error) {

	if err := validateStrategy(name, instruction, filters); err != nil {
		return nil, err
	}

	return &Strategy{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Instruction: instruction,
		Filters:     filters,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type Instruction struct {
	Selectors    map[string]string `json:"selectors"`
	OutputFormat map[string]string `json:"output_format"`
}

var exampleStrategyInstruction = &Instruction{
	Selectors: map[string]string{
		"#productTitle":        "text",
		"#priceblock_ourprice": "text",
	},

	OutputFormat: map[string]string{
		"#productTitle":        "string",
		"#priceblock_ourprice": "string",
	},
}

var (
	ErrNotFound = errors.New("extraction not found")
)

const (
	ExactStringMatch    StringMatchType = "exact"
	ContainsStringMatch StringMatchType = "contains"
)

const (
	ExactStringFilterType    FilterType = "string_exact"
	ContainsStringFilterType FilterType = "string_contains"
)

type Repository interface {
	Create(strategy *Strategy) error
	Get(id uuid.UUID) (*Strategy, error)
	ListByUserID(userID uuid.UUID) ([]*Strategy, error)
	Update(strategy *Strategy) error
}

type Service interface {
	Create(userID uuid.UUID, name string, selector string, filters []*Filter) (*Strategy, error)
	Get(id uuid.UUID) (*Strategy, error)
	ListByUserID(userID uuid.UUID) ([]*Strategy, error)
}

type Policy interface {
	CanCreate() error
	CanUpdate(strategy *Strategy) error
	CanRead(strategy *Strategy) error
	CanDelete(strategy *Strategy) error
}

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
}
