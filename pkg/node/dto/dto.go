package dto

type SelectorKey string
type FieldKey string
type Selector string

type ExtractionRequest struct {
	Selectors map[FieldKey]Selector `json:"selectors"`
}
