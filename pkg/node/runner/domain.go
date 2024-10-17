package runner

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Execute(c *gin.Context)
	Titles(c *gin.Context)
}

type Service interface {
	Titles() ([]string, error)
	Execute(src string) ([]byte, error)
}
