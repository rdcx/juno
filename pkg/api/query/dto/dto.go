package dto

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
	ID                string      `json:"id"`
	QueryType         QueryType   `json:"type"`
	BasicQueryVersion string      `json:"basic_query_version"`
	BasicQuery        *BasicQuery `json:"basic_query"`
}

type QueryResult interface{}

type GetResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Query Query `json:"query"`
}
