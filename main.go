package main

import (
	"alice/keramico/database"
	"alice/keramico/internal/redis"
	"alice/keramico/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	rdb := redis.NewRedisClient("localhost:6379", "", 0)
	defer rdb.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	routes.SetupRoutes(r, database.DB, rdb)

	r.Run(":8080")
}
