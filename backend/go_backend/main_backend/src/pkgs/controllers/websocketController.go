package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kimdwan/logan_drive/src/dtos"
	"github.com/kimdwan/logan_drive/src/pkgs/services"
)

type WebsocketController interface {
	WebsocketAuthFriendStatusController(ctx *gin.Context)
}

// 유저 친구들의 실시간을 확인할 수 있는 로직
func WebsocketAuthFriendStatusController(ctx *gin.Context) {

	var (
		conn          *websocket.Conn
		friend_status []dtos.WebsocketFriendDto
		errorStatus   int
		err           error
	)

	// 웹소켓 연결
	if conn, err = services.WebsocketConnectFunc(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer conn.Close()

	log.Println("웹소켓이 연결되었습니다")

	// 타이머 설정
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {

		friend_status = nil

		// 유저의 정보를 가져오는 로직
		if errorStatus, err = services.WebsocketAuthFriendStatusService(conn, &friend_status); err != nil {
			ctx.AbortWithStatusJSON(errorStatus, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 데이터를 보내는 로직
		if err = services.WebsocketSendDataService[[]dtos.WebsocketFriendDto](conn, &friend_status, websocket.TextMessage); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 시간 텀을 주는 로직
		time.Sleep(5 * time.Second)

	}

}
