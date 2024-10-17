package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type ExecuteRequest struct {
	Src string `json:"src"`
}

type ExecuteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Data []byte `json:"data,omitempty"`
}

func NewSuccessExecuteResponse(data []byte) ExecuteResponse {
	return ExecuteResponse{
		Status: SUCCESS,
		Data:   data,
	}
}

func NewErrorExecuteResponse(message string) ExecuteResponse {
	return ExecuteResponse{
		Status:  ERROR,
		Message: message,
	}
}
