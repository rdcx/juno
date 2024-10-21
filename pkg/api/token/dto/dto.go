package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type BalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Balance float64 `json:"balance"`
}

func NewSuccessBalanceResponse(balance float64) *BalanceResponse {
	return &BalanceResponse{
		Status:  SUCCESS,
		Balance: balance,
	}
}

func NewErrorBalanceResponse(message string) *BalanceResponse {
	return &BalanceResponse{
		Status:  ERROR,
		Message: message,
	}
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type DepositResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSuccessDepositResponse() *DepositResponse {
	return &DepositResponse{
		Status: SUCCESS,
	}
}

func NewErrorDepositResponse(message string) *DepositResponse {
	return &DepositResponse{
		Status:  ERROR,
		Message: message,
	}
}
