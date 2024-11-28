package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kimdwan/logan_drive/src/dtos"
	"github.com/kimdwan/logan_drive/src/pkgs/services"
)

type WebsocketController interface {
	WebsocketTestController(ctx *gin.Context)
	WebsocketUserStatusController(ctx *gin.Context)
}

// 테스트용 컨트롤러
func WebsocketTestController(ctx *gin.Context) {

	var (
		conn *websocket.Conn
		err  error
	)

	// 연결
	if conn, err = services.WebsocketTranslateService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// 초반
	if err = conn.WriteMessage(websocket.TextMessage, []byte(`연결되었습니다`)); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			if _, _, err = conn.ReadMessage(); err != nil {
				log.Println("시스템 오류: ", err.Error())
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			if err = conn.WriteMessage(websocket.TextMessage, []byte(`연결되었습니다`)); err != nil {
				log.Println("시스템 오류: ", err.Error())
				return
			}
		case <-c.Done():
			return
		}
	}

}

// 친구가 실시간으로 접속해 있는지 확인하는 함수
func WebsocketUserStatusController(ctx *gin.Context) {

	var (
		conn        *websocket.Conn
		errorStatus int
		err         error
	)

	// conn 연결
	if conn, err = services.WebsocketTranslateService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// 웹소켓 연결 실시간 확인
	var (
		user_computer_number *dtos.WebsocketUserComputerNumberDto
		friend_statuses      []dtos.WebsocketFriendStatusDto
		limit_count          int = 0
	)
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	// 클라이언트에 데이터를 실시간으로 읽고 전달하고 검정한다
	go func() {
		for {
			if user_computer_number, err = services.WebsocketParseDataService[dtos.WebsocketUserComputerNumberDto](conn, websocket.TextMessage); err != nil {
				if err = services.WebsocketTransformDataAndSendDataToClientService[dtos.WebsocketErrorPackDto](conn, &dtos.WebsocketErrorPackDto{Error: err.Error(), Status: http.StatusBadRequest}, websocket.TextMessage); err != nil {
					cancel()
					return
				}
				cancel()
				return
			} else {
				// 첫번째 서치
				if errorStatus, err = services.WebsocketUserStatusService(user_computer_number, &friend_statuses, &limit_count); err != nil {
					if err = services.WebsocketTransformDataAndSendDataToClientService[dtos.WebsocketErrorPackDto](conn, &dtos.WebsocketErrorPackDto{Error: err.Error(), Status: errorStatus}, websocket.TextMessage); err != nil {
						return
					}
					return
				} else {
					if err = services.WebsocketTransformDataAndSendDataToClientService[[]dtos.WebsocketFriendStatusDto](conn, &friend_statuses, websocket.TextMessage); err != nil {
						return
					}
				}
			}
		}
	}()

	// 시간 텀을 두고 메세지 전송
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if errorStatus, err = services.WebsocketUserStatusService(user_computer_number, &friend_statuses, &limit_count); err != nil {
				if err = services.WebsocketTransformDataAndSendDataToClientService[dtos.WebsocketErrorPackDto](conn, &dtos.WebsocketErrorPackDto{Error: err.Error(), Status: errorStatus}, websocket.TextMessage); err != nil {
					return
				}
				return
			} else {
				if err = services.WebsocketTransformDataAndSendDataToClientService[[]dtos.WebsocketFriendStatusDto](conn, &friend_statuses, websocket.TextMessage); err != nil {
					return
				}
			}
		case <-c.Done():
			return
		}

	}

}
