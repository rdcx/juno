package user

import (
	"context"
	"errors"
	"juno/pkg/can"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("user not found")
var ErrForbidden = errors.New("forbidden")
var ErrUserNotFoundInContext = errors.New("user not found in context")
var ErrInternal = errors.New("internal error")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidID = errors.New("invalid id")
var ErrInvalidName = errors.New("invalid name")

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type Repository interface {
	Create(u *User) error
	Get(id uuid.UUID) (*User, error)
	FirstWhereEmail(email string) (*User, error)
	Update(u *User) error
	Delete(id uuid.UUID) error
}

type Service interface {
	Create(name, email, password string) (*User, error)
	Get(id uuid.UUID) (*User, error)
	FirstWhereEmail(email string) (*User, error)
	Update(u *User) error
	Delete(id uuid.UUID) error
}

type Policy interface {
	CanCreate() can.Result
	CanUpdate(ctx context.Context, u *User) can.Result
	CanRead(ctx context.Context, u *User) can.Result
	CanDelete(ctx context.Context, u *User) can.Result
}

type Handler interface {
	Get(c *gin.Context)
	Profile(c *gin.Context)
	Create(c *gin.Context)
}
