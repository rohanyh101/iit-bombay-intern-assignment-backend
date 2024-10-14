package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/roh4nyh/iit_bombay/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controllers.UserSignUp())
	incomingRoutes.POST("users/login", controllers.UserLogIn())
}
