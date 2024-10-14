package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("expired token")

type Handler interface {
	Token(c *gin.Context)
}

type Service interface {
	Authenticate(email, password string) (string, error)
}
