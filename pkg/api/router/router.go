package router

import (
	"juno/pkg/api/node"
	"juno/pkg/api/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var UserID = "51318544-492d-4aa8-b474-ebc39b4c4773"

func New(
	nodeHandler node.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user", &user.User{
			ID: uuid.MustParse(UserID),
		})

		c.Next()
	})

	r.GET("/nodes/:id", nodeHandler.Get)
	r.POST("/nodes", nodeHandler.Create)
	r.PUT("/nodes", nodeHandler.Update)
	r.DELETE("/nodes/:id", nodeHandler.Delete)

	return r
}
