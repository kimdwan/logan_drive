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

	// 유저가 정보를 이용하는 라우터
	authuserrouter := authrouter.Group("user")
	authuserrouter.GET("logout", controllers.AuthUserLogoutController)
	authuserrouter.POST("upload/profile", controllers.AuthUserUploadProfileController)
	authuserrouter.GET("get/friendlist", controllers.AuthUserGetFriendListController)

	// 기본정보를 가져오는 라우터
	authgetrouter := authrouter.Group("get")
	authgetrouter.GET("detail", controllers.AuthGetUserEmailAndNickNameController)
	authgetrouter.GET("profileimg", controllers.AuthGetUserProfileImgController)

	// 친구창 관리와 관련된 라우터
	authfriendrouter := authrouter.Group("friend")
	authfriendrouter.POST("send/message", controllers.AuthFriendSendMessageController)
	authfriendrouter.POST("request", controllers.AuthFriendRequestController)
	authfriendrouter.POST("check/verify", controllers.AuthFriendCheckVerifyController)
}
