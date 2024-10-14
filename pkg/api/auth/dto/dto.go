package dto

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type TokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Token string `json:"token,omitempty"`
}

func NewSuccessTokenResponse(token string) TokenResponse {
	return TokenResponse{
		Status: SUCCESS,
		Token:  token,
	}
}

func NewErrorTokenResponse(message string) TokenResponse {
	return TokenResponse{
		Status:  ERROR,
		Message: message,
	}
}
