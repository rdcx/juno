package ranag

import (
	"juno/pkg/ranag/dto"

	"github.com/gin-gonic/gin"
)

type Service interface {
	RangeAggregate(offset, total int, req dto.RangeAggregatorRequest) ([]map[string]interface{}, error)
}

type Handler interface {
	RangeAggregate(c *gin.Context)
}
