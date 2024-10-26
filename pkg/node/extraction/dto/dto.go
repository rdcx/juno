package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Selector struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Field struct {
	ID         string `json:"id"`
	SelectorID string `json:"selector_id"`
	Name       string `json:"name"`
}

type ExtractionRequest struct {
	Selectors []*Selector `json:"selectors" binding:"required"`
	Fields    []*Field    `json:"fields" binding:"required"`
}

type ExtractionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Extractions []map[string]interface{} `json:"extractions,omitempty"`
}

func NewSuccessExtractionResponse(extractions []map[string]interface{}) *ExtractionResponse {
	return &ExtractionResponse{
		Status:      SUCCESS,
		Extractions: extractions,
	}
}

func NewErrorExtractionResponse(err error) *ExtractionResponse {
	return &ExtractionResponse{
		Status:  ERROR,
		Message: err.Error(),
	}
}
