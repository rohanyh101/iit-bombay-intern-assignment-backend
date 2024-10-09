package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/roh4nyh/iit_bombay/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	app := gin.New()
	app.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"Access-Control-Allow-Origin", "*"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	app.Use(cors.New(config))

	routes.AuthRoutes(app)

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "iit bombay server is up and running..."})
	})

	routes.LibrarianRoutes(app)

	routes.MemberRoutes(app)

	app.Run(fmt.Sprintf(":%s", PORT))
}
