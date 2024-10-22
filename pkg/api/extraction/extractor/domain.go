package extractor

import (
	"errors"
	"time"

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

type Extractor struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Selector  string    `json:"selector"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Create(extractor *Extractor) error
	Get(id uuid.UUID) (*Extractor, error)
	ListByUserID(userID uuid.UUID) ([]*Extractor, error)
	Update(extractor *Extractor) error
}

type Service interface {
	Create(userID uuid.UUID, name string, selector string, filters []*Filter) (*Extractor, error)
	Get(id uuid.UUID) (*Extractor, error)
	ListByUserID(userID uuid.UUID) ([]*Extractor, error)
}
