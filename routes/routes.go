package routes

import (
	"alice/keramico/handlers"
	"alice/keramico/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {

	r.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})
	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})

	r.GET("/products", func(c *gin.Context) {
		handlers.GetProducts(c, db)
	})

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		r.POST("/products", func(c *gin.Context) {
			handlers.CreateProduct(c, db)
		})
	}

}
