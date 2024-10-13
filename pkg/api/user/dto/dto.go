package dto

import "juno/pkg/api/user"

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type GetUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	User *User `json:"user,omitempty"`
}

func NewUserFromDomain(u *user.User) *User {
	return &User{
		ID:    u.ID.String(),
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
