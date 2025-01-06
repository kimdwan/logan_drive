package controllers

import (
	"context"
	"encoding/json"
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
	WebsocketFriendCheckMessagesController(ctx *gin.Context)
	WebsocketFriendAdmitFriendAppealController(ctx *gin.Context)
	WebsocketFriendConfirmPrivateController(ctx *gin.Context)
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

// 메세지를 실시간으로 확인하는 로직
func WebsocketFriendCheckMessagesController(ctx *gin.Context) {

	var (
		conn        *websocket.Conn
		errorStatus int
		err         error
	)

	// conn 가져오기
	if conn, err = services.WebsocketTranslateService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// 첫데이터 읽고 사용하기
	var (
		client_data   dtos.WebsocketFriendCheckDto
		message_datas []dtos.WebsocketFriendMessageDto
		errorMsg      dtos.WebsocketErrorPackDto
	)
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	// 첫번째로 보내야 하는 데이터
	go func() {
		for {
			// 클라이언트에서 보낸 데이터 읽기
			if err = client_data.WebsocketSendFriendMessageParseBodyFunc(conn); err != nil {
				log.Println(err.Error())
				cancel()
				return
			} else {
				// 첫번째 서치
				if errorStatus, err = services.WebsocketFriendCheckMessagesService(&client_data, &message_datas); err != nil {
					log.Println(err.Error())
					errorMsgByte, _ := json.Marshal(errorMsg)
					if err = conn.WriteMessage(websocket.TextMessage, errorMsgByte); err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					}
					cancel()
					return
				} else {
					msgDataByte, _ := json.Marshal(message_datas)
					if err = conn.WriteMessage(websocket.TextMessage, msgDataByte); err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					}
				}
			}
		}
	}()

	// 5초 간격 서치
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if errorStatus, err = services.WebsocketFriendCheckMessagesService(&client_data, &message_datas); err != nil {
				log.Println("시스템 오류: ", err.Error())
				errorMsg.Error = err.Error()
				errorMsg.Status = errorStatus
				errorMsgByte, _ := json.Marshal(errorMsg)
				if err = conn.WriteMessage(websocket.TextMessage, errorMsgByte); err != nil {
					log.Println("시스템 오류: ", err.Error())
					cancel()
					return
				}
				cancel()
				return
			} else {
				msg_data_byte, _ := json.Marshal(message_datas)
				if err = conn.WriteMessage(websocket.TextMessage, msg_data_byte); err != nil {
					log.Println("시스템 오류: ", err.Error())
					cancel()
					return
				}
			}
		case <-c.Done():
			return
		}
	}
}

// 친구창온거 실시간으로 확인
func WebsocketFriendAdmitFriendAppealController(ctx *gin.Context) {

	var (
		conn *websocket.Conn
		err  error
	)

	// conn 연결
	if conn, err = services.WebsocketTranslateService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// 보낸 데이터 확인
	var (
		computer_number_dto *dtos.WebsocketUserComputerNumberDto
		user_datas          []dtos.WebsocketStreamFriendAllowStatusDto
		errorStatus         int
	)
	c, cancel := context.WithCancel(ctx)

	defer cancel()
	go func() {
		for {
			if computer_number_dto, err = services.WebsocketParseDataService[dtos.WebsocketUserComputerNumberDto](conn, websocket.TextMessage); err != nil {
				cancel()
				return
			} else {
				// 첫번째 작동
				if errorStatus, err = services.WebsocketFriendAdmitFriendAppealService(computer_number_dto, &user_datas); err != nil {
					services.WebsocketSendErrorMsgService(conn, errorStatus, err)
					cancel()
					return
				} else {
					user_data_byte, err := json.Marshal(&user_datas)
					if err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					} else {
						if err = conn.WriteMessage(websocket.TextMessage, user_data_byte); err != nil {
							log.Println("시스템 오류: ", err.Error())
							cancel()
							return
						}
					}
				}
			}
		}
	}()

	// 주기적으로 데이터를 보내주는 함수
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if errorStatus, err = services.WebsocketFriendAdmitFriendAppealService(computer_number_dto, &user_datas); err != nil {
				services.WebsocketSendErrorMsgService(conn, errorStatus, err)
				cancel()
				return
			} else {
				user_data_byte, err := json.Marshal(&user_datas)
				if err != nil {
					log.Println("시스템 오류: ", err.Error())
					cancel()
					return
				} else {
					if err = conn.WriteMessage(websocket.TextMessage, user_data_byte); err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					}
				}
			}
		case <-c.Done():
			return
		}
	}
}

// 유저 한명의 정보를 가져오는 로직
func WebsocketFriendConfirmPrivateController(ctx *gin.Context) {

	var (
		conn *websocket.Conn
		err  error
	)

	// conn 가져오기
	if conn, err = services.WebsocketTranslateService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer conn.Close()

	// 초반에 데이터 가져오기
	var (
		computerNumberAndFriendId dtos.WebsocketComputerNumberAndFriendIdDto
		friend_detail_datas       dtos.WebsocketCheckFriendDetailDto
		errorStatus               int
	)
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	// 초기 문제
	go func() {
		for {
			if err = computerNumberAndFriendId.WebsocketComputerNumberAndFriendIdParseDataAndCheckValidateAndSearchFunc(conn); err != nil {
				services.WebsocketSendErrorMsgService(conn, http.StatusBadRequest, err)
				cancel()
				return
			} else {
				if errorStatus, err = services.WebsocketFriendConfirmPrivateService(&computerNumberAndFriendId, &friend_detail_datas); err != nil {
					services.WebsocketSendErrorMsgService(conn, errorStatus, err)
					cancel()
					return
				} else {
					if data_byte, err := json.Marshal(&friend_detail_datas); err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					} else {
						if err = conn.WriteMessage(websocket.TextMessage, data_byte); err != nil {
							log.Println("시스템 오류: ", err.Error())
							cancel()
							return
						}
					}
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if errorStatus, err = services.WebsocketFriendConfirmPrivateService(&computerNumberAndFriendId, &friend_detail_datas); err != nil {
				services.WebsocketSendErrorMsgService(conn, errorStatus, err)
				cancel()
				return
			} else {
				if data_byte, err := json.Marshal(&friend_detail_datas); err != nil {
					log.Println("시스템 오류: ", err.Error())
					cancel()
					return
				} else {
					if err = conn.WriteMessage(websocket.TextMessage, data_byte); err != nil {
						log.Println("시스템 오류: ", err.Error())
						cancel()
						return
					}
				}
			}
		case <-c.Done():
			return
		}
	}

}
