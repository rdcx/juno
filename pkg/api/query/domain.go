package query

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrQueryNotFound = errors.New("query not found")
)

type StringMatchType string

const (
	ExactStringMatch    StringMatchType = "exact"
	ContainsStringMatch StringMatchType = "contains"
)

type StringMatch struct {
	Value     string          `json:"value"`
	MatchType StringMatchType `json:"type"`
}

type LinkMatch struct {
	Src *StringMatch `json:"src"`
	Dst *StringMatch `json:"dst"`
}

type BasicQuery struct {
	Title       *StringMatch `json:"title"`
	Description *StringMatch `json:"description"`
	Links       []*LinkMatch `json:"links"`
}

type Status string

const (
	PendingStatus   Status = "pending"
	RunningStatus   Status = "running"
	CompletedStatus Status = "completed"
	FailedStatus    Status = "failed"
)

type QueryType string

const (
	BasicQueryType QueryType = "basic"
)

type Query struct {
	ID                uuid.UUID   `json:"id"`
	UserID            uuid.UUID   `json:"user_id"`
	Status            Status      `json:"status"`
	QueryType         QueryType   `json:"type"`
	BasicQueryVersion string      `json:"basic_query_version"`
	BasicQuery        *BasicQuery `json:"basic_query"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(query *Query) error
	Get(id uuid.UUID) (*Query, error)
	ListByUserID(userID uuid.UUID) ([]*Query, error)
	Update(query *Query) error
}

type Service interface {
	Create(userID uuid.UUID, basicQuery *BasicQuery) (*Query, error)
	Get(id uuid.UUID) (*Query, error)
	ListByUserID(userID uuid.UUID) ([]*Query, error)
}
