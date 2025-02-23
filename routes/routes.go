package routes

import (
	"alice/keramico/handlers"
	"alice/keramico/internal/redis"
	"alice/keramico/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, redisClient *redis.RedisClient) {

	r.POST("/api/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})
	r.POST("/api/login", func(c *gin.Context) {
		handlers.Login(c, db, redisClient)
	})

	r.GET("/api/products", func(c *gin.Context) {
		handlers.GetProducts(c, db)
	})

	r.GET("/api/products/:id", func(c *gin.Context) {
		handlers.GetProductById(c, db)
	})

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(redisClient))
	{
		auth.POST("/api/products", func(c *gin.Context) {
			handlers.CreateProduct(c, db)
		})
		auth.POST("/api/logout", func(c *gin.Context) {
			handlers.Logout(c, redisClient)
		})
	}

}
