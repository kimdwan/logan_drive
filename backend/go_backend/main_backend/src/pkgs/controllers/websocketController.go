package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kimdwan/logan_drive/src/dtos"
	"github.com/kimdwan/logan_drive/src/pkgs/services"
)

type WebsocketController interface {
	WebsocketTestController(ctx *gin.Context)
	WebsocketAuthFriendStatusController(ctx *gin.Context)
}

// 테스트용 로직
func WebsocketTestController(ctx *gin.Context) {

	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			// origin header 파싱
			origin := r.Header.Get("Origin")

			parse_url, err := url.Parse(origin)
			if err != nil {
				log.Println("시스템 오류: ", err.Error())
				return false
			}

			url_name := parse_url.Hostname()

			// 검증
			var (
				allowed_hosts []string = strings.Split(os.Getenv("GO_ALLOWED_HOST_NAME"), ",")
			)
			for _, allowed_host := range allowed_hosts {
				if url_name == allowed_host {
					return true
				}
			}

			return false
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "웹소켓을 연결하는데 오류가 발생했습니다",
		})
		return
	}

	defer conn.Close()
	ticker := time.NewTicker(5 * time.Second)
	var stopConnect = make(chan bool)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.WriteMessage(websocket.TextMessage, []byte(`연결 되었습니다`))
		case <-stopConnect:
			return
		}
	}

}

// 유저 친구들의 실시간을 확인할 수 있는 로직
func WebsocketAuthFriendStatusController(ctx *gin.Context) {

	var (
		conn        *websocket.Conn
		errorStatus int
		err         error
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

	for {

		var friend_status []dtos.WebsocketFriendDto

		// 유저의 정보를 가져오는 로직
		if errorStatus, err = services.WebsocketAuthFriendStatusService(conn, &friend_status); err != nil {
			errorMsg, _ := json.Marshal(map[string]string{
				"error":  err.Error(),
				"status": strconv.Itoa(errorStatus),
			})
			conn.WriteMessage(websocket.TextMessage, errorMsg)
			return
		}

		// 데이터를 보내는 로직
		if err = services.WebsocketSendDataService[[]dtos.WebsocketFriendDto](conn, &friend_status, websocket.TextMessage); err != nil {
			errorMsg, _ := json.Marshal(map[string]string{
				"error":  err.Error(),
				"status": strconv.Itoa(http.StatusInternalServerError),
			})
			conn.WriteMessage(websocket.TextMessage, errorMsg)
			return
		}

	}

}
