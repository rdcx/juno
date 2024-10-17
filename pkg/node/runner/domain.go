package runner

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Execute(c *gin.Context)
}

type Service interface {
	Execute(src string) ([]byte, error)
}
