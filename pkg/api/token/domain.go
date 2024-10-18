package token

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrInvalidAmount = errors.New("invalid amount")

type Service interface {
	Balance(userID uuid.UUID) (int, error)
	Deposit(userID uuid.UUID, amount int) error
}

type Handler interface {
	Balance(c *gin.Context)
	Deposit(c *gin.Context)
}
