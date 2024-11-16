package info

import "github.com/gin-gonic/gin"

type Info struct {
	PageCount int `json:"page_count"`
}

type Handler interface {
	Info(c *gin.Context)
}

type Service interface {
	GetInfo() (*Info, error)
}
