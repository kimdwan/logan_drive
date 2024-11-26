package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/pkgs/controllers"
)

func WebsocketRouter(router *gin.Engine) {
	websocketrouter := router.Group("ws")

	// 테스트용
	websocketrouter.GET("test", controllers.WebsocketTestController)

	// 유저가 사용하는 라우터
	wsuserrouter := websocketrouter.Group("user")
	wsuserrouter.GET("friends/connect", controllers.WebsocketAuthFriendStatusController)

}
