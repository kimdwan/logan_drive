package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/dtos"
	"gorm.io/gorm"
)

type WebsocketService interface {
	WebsocketConnectFunc(ctx *gin.Context) error
	WebsocketReadDataService(conn *websocket.Conn, data *[]byte, dataType int) error
	WebsocketAuthFriendStatusService(conn *websocket.Conn, friend_status *[]dtos.WebsocketFriendDto) (int, error)
}

// 웹소켓 연결을 진행해 주는 함수
func WebsocketConnectFunc(ctx *gin.Context) (*websocket.Conn, error) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {

			// 클라이언트에 url을 가져와서 확인하기
			origin := r.Header.Get("Origin")

			parse_url, err := url.Parse(origin)
			if err != nil {
				log.Println("시스템 오류: ", err.Error())
				return false
			}

			url_name := parse_url.Hostname()
			log.Println(url_name)

			// header 검증하기
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
		return nil, errors.New("websocket연결을 하는데 오류가 발생했습니다")
	}

	return conn, nil
}

// 기본적으로 데이터를 읽는 함수
func WebsocketReadDataService(conn *websocket.Conn, data *[]byte, dataType int) error {
	var (
		wsDataType int
		err        error
	)

	// conn에서 보낸 데이터 확인
	wsDataType, *data, err = conn.ReadMessage()
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("웹소켓에서 보낸 데이터를 읽는데 오류가 발생했습니다")
	}

	// 데이터가 문자열이 맞는지 확인
	if wsDataType != dataType {
		log.Println("시스템 오류: 웹소켓 데이터 타입 오류")
		return errors.New("웹소켓에서 보낸 데이터 타입이 텍스트 타입이 아닙니다")
	}

	return nil

}

// 데이터 보내기
func WebsocketSendDataService[T []dtos.WebsocketFriendDto](conn *websocket.Conn, send_data *T, dataType int) error {

	var (
		datas []byte
		err   error
	)

	// 데이터 직렬화
	if datas, err = json.Marshal(send_data); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("데이터를 파싱하는데 오류가 발생했습니다")
	}

	// 데이터를 보내기
	if err = conn.WriteMessage(dataType, datas); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("데이터를 클라이언트에 보내는데 오류가 발생했습니다")
	}

	return nil
}

// 유저의 친구창을 확인 후 상태창을 보여주는 함수
func WebsocketAuthFriendStatusService(conn *websocket.Conn, friend_status *[]dtos.WebsocketFriendDto) (int, error) {

	var (
		wsData      []byte
		dataType    = websocket.TextMessage
		errorStatus int
		err         error
	)

	// 웹소켓 에서 데이터 가져오기
	if err = WebsocketReadDataService(conn, &wsData, dataType); err != nil {
		return http.StatusBadRequest, err
	}

	var (
		db           *gorm.DB = settings.DB
		friend_lists []uuid.UUID
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 웹소켓에 데이터를 파싱하고 친구창 데이터 가져오기
	if errorStatus, err = WebsocketAuthFreindStatusParseDataAndFindFriendFunc(c, db, &wsData, &friend_lists); err != nil {
		return errorStatus, err
	}

	// 친구가 없다면 스킵
	if len(friend_lists) == 0 {
		return 0, nil
	}

	// 본격적으로 친구에 접속여부를 확인하는 로직
	if err = WebsocketAuthFriendStatusGetDataFunc(c, db, &friend_lists, friend_status); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

type WebsocketAuthFriendStatus interface {
	WebsocketAuthFreindStatusParseDataAndFindFriendFunc(c context.Context, db *gorm.DB, wsData *[]byte, friend_lists *[]uuid.UUID) (int, error)
	WebsocketAuthFreindStatusParseFriendDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_data_lists *[]servicemodel.Friend, friend_lists *[]uuid.UUID)
	WebsocketAuthFriendStatusGetDataFunc(c context.Context, db *gorm.DB, friend_lists *[]uuid.UUID, friend_status *[]dtos.WebsocketFriendDto) error
	WebsocketAuthFriendStatusGetDataSyncFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, mutex *sync.Mutex, friend_lists *[]uuid.UUID, friend_status *[]dtos.WebsocketFriendDto, error_lists *[]error)
}

// 웹소켓에서 보낸 데이터를 파싱하고 친구창 데이터를 확인함
func WebsocketAuthFreindStatusParseDataAndFindFriendFunc(c context.Context, db *gorm.DB, wsData *[]byte, friend_lists *[]uuid.UUID) (int, error) {

	var computer_number dtos.WebsocketComputerNumberDto

	// wsData를 파싱해주어야 함
	err := json.Unmarshal(*wsData, &computer_number)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return http.StatusBadRequest, errors.New("컴퓨터 넘버를 파싱하는데 오류가 발생했습니다")
	}

	// 유저 정보부터 찾기
	var user servicemodel.User
	result := db.WithContext(c).Where("computer_number = ?", computer_number.Computer_number).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("시스템 오류: 데이터 베이스에서 클라이언트에서 보낸 computer number와 일치하는 유저 데이터를 찾지 못함")
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 computer number를 다시 확인하세요")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저 정보를 찾는데 오류가 발생했습니다")
		}
	}

	// 데이터 베이스에서 친구 정보를 찾기
	var friend_data_lists []servicemodel.Friend
	if result = db.WithContext(c).Where("friend_1 = ? OR friend_2 = ?", user.User_id, user.User_id).Find(&friend_data_lists); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return http.StatusInternalServerError, errors.New("데이터 베이스에서 친구에 정보를 찾는데 오류가 발생했습니다")
	}

	// 친구가 없으면 skip
	if len(friend_data_lists) == 0 {
		return 0, nil
	}

	// 친구 데이터 정리
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	wg.Add(1)
	go WebsocketAuthFreindStatusParseFriendDataFunc(&wg, &mutex, &user.User_id, &friend_data_lists, friend_lists)
	wg.Wait()

	return 0, nil

}

// 데이터 베이스에서 보낸 친구 데이터를 파싱하는데 사용
func WebsocketAuthFreindStatusParseFriendDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_data_lists *[]servicemodel.Friend, friend_lists *[]uuid.UUID) {

	defer wg.Done()
	for _, friend_data := range *friend_data_lists {

		mutex.Lock()
		if *user_id == friend_data.Friend_1 {
			*friend_lists = append(*friend_lists, friend_data.Friend_2)
		} else {
			*friend_lists = append(*friend_lists, friend_data.Friend_1)
		}
		mutex.Unlock()

	}

}

// 친구의 상태창을 확인하는데 사용
func WebsocketAuthFriendStatusGetDataFunc(c context.Context, db *gorm.DB, friend_lists *[]uuid.UUID, friend_status *[]dtos.WebsocketFriendDto) error {

	var (
		wg          sync.WaitGroup
		mutex       sync.Mutex
		error_lists []error
	)

	// 친구 데이터 데이터 베이스에서 찾고 가져오기
	wg.Add(1)
	go WebsocketAuthFriendStatusGetDataSyncFunc(c, db, &wg, &mutex, friend_lists, friend_status, &error_lists)
	wg.Wait()

	if len(error_lists) != 0 {
		for _, error_value := range error_lists {
			log.Println("시스템 오류: ", error_value.Error())
		}
		return errors.New("데이터 베이스에서 유저의 정보를 가져오는데 오류가 발생했습니다")
	}

	return nil
}

// 친구 상태창 확인 후 정리
func WebsocketAuthFriendStatusGetDataSyncFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, mutex *sync.Mutex, friend_lists *[]uuid.UUID, friend_status *[]dtos.WebsocketFriendDto, error_lists *[]error) {

	defer wg.Done()
	for _, friend := range *friend_lists {

		mutex.Lock()
		// 데이터 베이스에서 유저 정보 가져오기
		var user servicemodel.User
		result := db.WithContext(c).Where("user_id = ?", friend).First(&user)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				*error_lists = append(*error_lists, result.Error)
			} else {
				log.Println("존재하지 않는 유저 입니다")
			}
		}

		// 상태에 따른 값 배정 로그아웃 = 0, 활성화 = 1, 부재중(5분 정도) = 2, 연결안됨(1시간 동안 활동 안할때) = 3
		var status dtos.WebsocketFriendDto
		if user.Computer_number != nil {

			var (
				now = time.Now()
			)
			if now.Before(user.UpdatedAt.Add(time.Duration(5) * time.Minute)) {
				status.Friend_id = user.User_id
				status.Friend_status = 1
			} else {
				if now.Before(user.UpdatedAt.Add(time.Duration(1) * time.Hour)) {
					status.Friend_id = user.User_id
					status.Friend_status = 2
				} else {
					status.Friend_id = user.User_id
					status.Friend_status = 3
				}
			}

		} else {
			status.Friend_id = user.User_id
			status.Friend_status = 0
		}

		*friend_status = append(*friend_status, status)

		mutex.Unlock()
	}

}
