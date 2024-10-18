package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type BalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Balance int `json:"balance"`
}

func NewSuccessBalanceResponse(balance int) *BalanceResponse {
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
	Amount int `json:"amount"`
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
