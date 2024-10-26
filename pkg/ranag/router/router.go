package router

import (
	"juno/pkg/ranag"

	"github.com/gin-gonic/gin"
)

func New(
	handler ranag.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/aggregate", handler.RangeAggregate)

	return r
}
