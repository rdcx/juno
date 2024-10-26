package extraction

import (
	"github.com/gin-gonic/gin"

	"juno/pkg/node/extraction/dto"
)

type Handler interface {
	Extract(c *gin.Context)
}

type Service interface {
	Extract(req dto.ExtractionRequest) ([]map[string]interface{}, error)
}
