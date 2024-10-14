package router

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/balancer"
	"juno/pkg/api/middleware"
	"juno/pkg/api/node"
	"juno/pkg/api/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func New(
	nodeHandler node.Handler,
	userHandler user.Handler,
	authHandler auth.Handler,
	balancerHandler balancer.Handler,
) *gin.Engine {
	r := gin.Default()

	godotenv.Load("api.env")

	r.POST("/auth/token", authHandler.Token)
	r.POST("/users", userHandler.Create)

	authGroup := r.Group("/")

	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/users/:id", userHandler.Get)

		authGroup.GET("/nodes/:id", nodeHandler.Get)
		authGroup.POST("/nodes", nodeHandler.Create)
		authGroup.PUT("/nodes", nodeHandler.Update)
		authGroup.DELETE("/nodes/:id", nodeHandler.Delete)

		authGroup.GET("/balancers/:id", balancerHandler.Get)
		authGroup.POST("/balancers", balancerHandler.Create)
		authGroup.PUT("/balancers", balancerHandler.Update)
		authGroup.DELETE("/balancers/:id", balancerHandler.Delete)
	}

	return r
}
