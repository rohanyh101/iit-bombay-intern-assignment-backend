package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/roh4nyh/iit_bombay/controllers"
	"github.com/roh4nyh/iit_bombay/middleware"
)

func LibrarianRoutes(incomingRoutes *gin.Engine) {
	librarianRoutes := incomingRoutes.Group("/librarian")
	librarianRoutes.Use(middleware.Authenticate(), middleware.AuthenticateLibrarian())

	// librarian CRUD operations
	librarianRoutes.POST("/books", controller.AddBook())
	librarianRoutes.GET("/books", controller.GetBooks())
	librarianRoutes.GET("/books/:isbn", controller.GetBook())
	librarianRoutes.PUT("/books/:isbn", controller.UpdateBook())
	librarianRoutes.DELETE("/books/:isbn", controller.DeleteBook())

	// member CRUD operations
	librarianRoutes.GET("/users", controller.GetUsers())
	librarianRoutes.POST("/users", controller.AddUser())
	librarianRoutes.GET("/users/:user_id", controller.GetUser())
	librarianRoutes.PUT("/users/:user_id", controller.UpdateUser())
	librarianRoutes.DELETE("/users/:user_id", controller.DeActivateUser())
	// force delete user (optional)
	librarianRoutes.DELETE("/users/:user_id/force", controller.DeleteUser())

	// get active users
	librarianRoutes.GET("/users/active", controller.GetActiveUsers())

	// get deleted users
	librarianRoutes.GET("/users/deleted", controller.GetNonActiveUsers())

	// member borrowed history
	librarianRoutes.GET("/users/:user_id/history", controller.GetTransactionHistory())
}
