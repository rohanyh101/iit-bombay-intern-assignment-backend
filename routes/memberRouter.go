package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/roh4nyh/iit_bombay/controllers"
	"github.com/roh4nyh/iit_bombay/middleware"
)

func MemberRoutes(incomingRoutes *gin.Engine) {
	memberRoutes := incomingRoutes.Group("/member")
	memberRoutes.Use(middleware.Authenticate(), middleware.AuthenticateMember())

	// user crud
	memberRoutes.GET("/books", controller.GetBooks())
	memberRoutes.GET("/books/:isbn", controller.GetBook())

	// member crud
	memberRoutes.POST("/books/borrow/:isbn", controller.BorrowBook())
	memberRoutes.PUT("/books/return/:isbn", controller.ReturnBook())
	memberRoutes.GET("/books/borrowed", controller.BorrowedBooks())
	memberRoutes.DELETE("/account", controller.DeActivateMember())
}
