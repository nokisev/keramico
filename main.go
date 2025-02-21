package main

import (
	"alice/keramico/database"
	"alice/keramico/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()
	defer database.DB.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	routes.SetupRoutes(r, database.DB)

	r.Run(":8080")
}
