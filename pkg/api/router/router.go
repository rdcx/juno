package router

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/balancer"
	"juno/pkg/api/middleware"
	"juno/pkg/api/node"
	"juno/pkg/api/token"
	"juno/pkg/api/transaction"
	"juno/pkg/api/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func New(
	nodeHandler node.Handler,
	balancerHandler balancer.Handler,
	transactionHandler transaction.Handler,
	tokenHandler token.Handler,
	userHandler user.Handler,
	authHandler auth.Handler,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware configuration
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,                                                                  // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},          // Allow all methods
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // Allow common headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Preflight request cache duration
	}))

	godotenv.Load("api.env")

	r.POST("/auth/token", authHandler.Token)
	r.POST("/users", userHandler.Create)
	r.GET("/shards/nodes", nodeHandler.AllShardsNodes)
	r.GET("/shards/balancers", balancerHandler.AllShardsBalancers)

	authGroup := r.Group("/")

	authGroup.Use(middleware.AuthMiddleware())
	{

		authGroup.GET("/profile", userHandler.Profile)

		authGroup.GET("/users/:id", userHandler.Get)

		authGroup.GET("/nodes", nodeHandler.List)
		authGroup.GET("/nodes/:id", nodeHandler.Get)
		authGroup.POST("/nodes", nodeHandler.Create)
		authGroup.PUT("/nodes/:id", nodeHandler.Update)
		authGroup.DELETE("/nodes/:id", nodeHandler.Delete)

		authGroup.GET("/balancers", balancerHandler.List)
		authGroup.GET("/balancers/:id", balancerHandler.Get)
		authGroup.POST("/balancers", balancerHandler.Create)
		authGroup.PUT("/balancers/:id", balancerHandler.Update)

		authGroup.GET("/tokens/balance", tokenHandler.Balance)
		authGroup.POST("/tokens/deposit", tokenHandler.Deposit)

		authGroup.GET("/transactions", transactionHandler.List)
	}

	return r
}
