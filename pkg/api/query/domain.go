package query

import "github.com/google/uuid"

type StringMatchType string

const (
	ExactStringMatch    StringMatchType = "exact"
	ContainsStringMatch StringMatchType = "contains"
)

type StringMatch struct {
	Value string          `json:"value"`
	Type  StringMatchType `json:"type"`
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

type Query struct {
	UserID uuid.UUID `json:"user_id"`

	BasicQuery *BasicQuery
}
