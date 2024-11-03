package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/middlewares"
	"github.com/kimdwan/logan_drive/src/pkgs/controllers"
)

// 로그인한 유저만 사용할 수 있는 로직
func AuthRouter(router *gin.Engine) {

	authrouter := router.Group("auth")
	authrouter.Use(middlewares.CheckJwtMiddleware())

	// 기본정보를 가져오는 라우터
	authgetrouter := authrouter.Group("get")
	authgetrouter.GET("detail", controllers.AuthGetUserEmailAndNickNameController)
	authgetrouter.GET("profileimg", controllers.AuthGetUserProfileImgController)
}
