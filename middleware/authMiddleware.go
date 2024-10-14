package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/roh4nyh/iit_bombay/helpers"
	"github.com/roh4nyh/iit_bombay/models"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the Authorization header instead of token
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header found"})
			c.Abort()
			return
		}

		// Remove "Bearer " from the token string
		token := clientToken[len("Bearer "):]

		claims, err := helper.ValidateUserToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)
		c.Set("uid", claims.Uid)
		c.Set("is_active", claims.IsActive)

		log.Printf("%+v", claims)

		c.Next()
	}
}

func AuthenticateLibrarian() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != models.ROLE_LIBRARIAN {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthenticated to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthenticateMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != models.ROLE_MEMBER {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthenticated to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
