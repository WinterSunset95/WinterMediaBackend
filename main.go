package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/WinterSunset95/WinterMediaBackend/api"
	"github.com/WinterSunset95/WinterMediaBackend/cognito"
	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	_ = port

	database.Init()
	cognito.InitCognito()
	certManager := autocert.Manager{
		Prompt:    autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("localhost", "ec2-13-235-70-83.ap-south-1.compute.amazonaws.com"),
		Cache: autocert.DirCache("certs"),
	}

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

	autotls.RunWithManager(app, &certManager)
}
