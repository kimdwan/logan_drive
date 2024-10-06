package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/pkgs/controllers"
)

func UserRouter(router *gin.Engine) {

	// 기본적인 유저에게 제공되는 라우터
	userrouter := router.Group("user")

	userrouter.POST("signup", controllers.UserSignUpController)
	userrouter.POST("login", controllers.UserLoginController)
}
