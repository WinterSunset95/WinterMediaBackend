package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/WinterSunset95/WinterMediaBackend/api"
	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")

	database.Init()
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowCredentials: true,
	}))

	api.ApiRoutes(app)

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	app.GET("/getmovies", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the getmovies endpoint",
		})
	})

	app.Run(":" + port)
}
