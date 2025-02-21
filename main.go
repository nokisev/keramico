package main

import (
	"alice/keramico/database"
	"alice/keramico/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	r := gin.Default()

	routes.SetupRoutes(r, database.DB)

	r.Run(":8080")
}
