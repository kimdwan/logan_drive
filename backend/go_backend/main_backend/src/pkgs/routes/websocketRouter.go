package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/pkgs/controllers"
)

func WebsocketRouter(router *gin.Engine) {

	wsrouter := router.Group("ws")

	// 테스트용
	wsrouter.GET("test", controllers.WebsocketTestController)

	// 유저가 사용하는 웹소켓
	wsuserrouter := wsrouter.Group("user")
	wsuserrouter.GET("status", controllers.WebsocketUserStatusController)
}
