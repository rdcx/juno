package dto

import "juno/pkg/api/transaction"

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Transaction struct {
	ID     string            `json:"id"`
	UserID string            `json:"user_id"`
	Amount float64           `json:"amount"`
	Key    string            `json:"key"`
	Meta   map[string]string `json:"meta"`
}

type ListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Transactions []Transaction `json:"transactions"`
}

func NewSuccessListResponse(transactions []*transaction.Transaction) *ListResponse {

	var dtoTransactions []Transaction
	for _, t := range transactions {
		dtoTransactions = append(dtoTransactions, Transaction{
			ID:     t.ID.String(),
			UserID: t.UserID.String(),
			Amount: t.Amount,
			Key:    string(t.Key),
			Meta:   t.Meta,
		})
	}

	return &ListResponse{
		Status:       SUCCESS,
		Transactions: dtoTransactions,
	}
}

func NewErrorListResponse(message string) *ListResponse {
	return &ListResponse{
		Status:  ERROR,
		Message: message,
	}
}
