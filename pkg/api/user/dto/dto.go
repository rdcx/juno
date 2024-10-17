package dto

import "juno/pkg/api/user"

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type User struct {
	ID    string `json:"id" validate:"required,uuid"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type GetUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	User *User `json:"user,omitempty"`
}

func NewUserFromDomain(u *user.User) *User {
	return &User{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
	}
}

func NewSuccessGetUserResponse(u *user.User) GetUserResponse {
	user := NewUserFromDomain(u)
	return GetUserResponse{
		Status: SUCCESS,
		User:   user,
	}
}

func NewErrorGetUserResponse(message string) GetUserResponse {
	return GetUserResponse{
		Status:  ERROR,
		Message: message,
	}
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	User *User `json:"user,omitempty"`
}

func NewSuccessCreateUserResponse(u *user.User) CreateUserResponse {
	user := NewUserFromDomain(u)
	return CreateUserResponse{
		Status: SUCCESS,
		User:   user,
	}
}

func NewErrorCreateUserResponse(message string) CreateUserResponse {
	return CreateUserResponse{
		Status:  ERROR,
		Message: message,
	}
}
