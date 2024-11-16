package dto

import "juno/pkg/node/info"

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type Info struct {
	PageCount int `json:"page_count"`
}

type InfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	Info *Info `json:"info,omitempty"`
}

func NewInfoFromDomain(info *info.Info) *Info {
	return &Info{
		PageCount: info.PageCount,
	}
}

func NewSuccessInfoResponse(info *info.Info) *InfoResponse {
	return &InfoResponse{
		Status: SUCCESS,
		Info:   NewInfoFromDomain(info),
	}
}

func NewErrorInfoResponse(message string) *InfoResponse {
	return &InfoResponse{
		Status:  ERROR,
		Message: message,
	}
}
