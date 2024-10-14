package router

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/middleware"
	"juno/pkg/api/node"
	"juno/pkg/api/user"

	"github.com/gin-gonic/gin"
)

func New(
	nodeHandler node.Handler,
	userHandler user.Handler,
	authHandler auth.Handler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/auth/token", authHandler.Token)
	r.POST("/users", userHandler.Create)

	authGroup := r.Group("/")

	authGroup.Use(middleware.AuthMiddleware())
	{
		r.GET("/nodes/:id", nodeHandler.Get)
		r.POST("/nodes", nodeHandler.Create)
		r.PUT("/nodes", nodeHandler.Update)
		r.DELETE("/nodes/:id", nodeHandler.Delete)

		r.GET("/users/:id", userHandler.Get)
	}

	return r
}
