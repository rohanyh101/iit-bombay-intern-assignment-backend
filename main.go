package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/roh4nyh/iit_bombay/middleware"
	"github.com/roh4nyh/iit_bombay/routes"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Printf("error loading .env file: %v", err)
	// }

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	app.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	app.Use(cors.New(config))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "iit bombay server is up and running..."})
	})

	routes.AuthRoutes(app)

	routes.LibrarianRoutes(app)

	routes.MemberRoutes(app)

	app.GET("/api/v1/whoami", middleware.Authenticate(), func(c *gin.Context) {
		username := c.GetString("username")
		role := c.GetString("role")
		c.JSON(http.StatusOK, gin.H{"username": username, "role": role})
	})

	app.Run(fmt.Sprintf(":%s", PORT))
}
