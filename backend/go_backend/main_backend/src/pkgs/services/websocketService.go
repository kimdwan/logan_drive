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
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/dtos"
	"gorm.io/gorm"
)

type WebsocketService interface {
	WebsocketTranslateService(ctx *gin.Context) (*websocket.Conn, error)
	WebsocketUserStatusService(user_computer_number *dtos.WebsocketUserComputerNumberDto, friend_statuses *[]dtos.WebsocketFriendStatusDto, limit_count int) (int, error)
}

// websokcet으로 변환 시켜주는 함수
func WebsocketTranslateService(ctx *gin.Context) (*websocket.Conn, error) {

	var (
		websocketUpgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     WebsocketTranslateCheckOriginFunc,
		}
	)

	// 본격적인 변환
	conn, err := websocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return nil, errors.New("웹소켓을 변환하는데 오류가 발생했습니다")
	}

	return conn, nil
}

type WebsocketTranslate interface {
	WebsocketTranslateCheckOriginFunc(r *http.Request) bool
}

// origin을 확인하는 함수
func WebsocketTranslateCheckOriginFunc(r *http.Request) bool {

	// origin header 추출
	origin := r.Header.Get("Origin")

	// origin 파싱
	parse_url, err := url.Parse(origin)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return false
	}

	origin_host := parse_url.Hostname()

	// 검증
	var (
		allowed_hosts []string = strings.Split(os.Getenv("GO_ALLOWED_HOST_NAME"), ",")
	)
	for _, allowed_host := range allowed_hosts {
		if origin_host == allowed_host {
			return true
		}
	}

	log.Println("시스템 오류: 허용되지 않는 host이름 입니다")
	return false
}

// 데이터를 읽는 함수
func WebsocketParseDataService[T dtos.WebsocketUserComputerNumberDto](conn *websocket.Conn, wantDataType int) (*T, error) {

	var (
		client_data T
		dataType    int
		data        []byte
		err         error
	)

	// 데이터 읽기
	if dataType, data, err = conn.ReadMessage(); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return nil, errors.New("데이터를 읽는데 오류가 발생했습니다")
	}

	// 데이터 검정
	if dataType != wantDataType {
		log.Printf("시스템 오류: 데이터 타입이 %v이 아님", wantDataType)
		return nil, errors.New("데이터 타입을 다시 확인해주세요")
	}

	// 데이터 변환
	if err = json.Unmarshal(data, &client_data); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return nil, errors.New("클라이언트에 데이터를 변환하는데 오류가 발생했습니다")
	}

	// 데이터 검정2
	validate := validator.New()
	if err = validate.Struct(client_data); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return nil, errors.New("클라이언트 데이터를 검정하는데 오류가 발생했습니다")
	}

	return &client_data, nil
}

// 웹소켓에서 실시간으로 데이터를 보내주는 함수
func WebsocketTransformDataAndSendDataToClientService[T dtos.WebsocketErrorPackDto | []dtos.WebsocketFriendStatusDto](conn *websocket.Conn, data *T, dataType int) error {

	// 데이터 변환
	data_byte, err := json.Marshal(data)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("데이터를 변환하는데 오류가 발생했습니다")
	}

	// 데이터 전송
	if err = conn.WriteMessage(dataType, data_byte); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("데이터를 전송하는데 오류가 발생했습니다")
	}

	return nil
}

// 유저의 친구들이 실시간으로 접속해 있는지 확인하는 함수
func WebsocketUserStatusService(user_computer_number *dtos.WebsocketUserComputerNumberDto, friend_statuses *[]dtos.WebsocketFriendStatusDto, limit_count *int) (int, error) {

	var (
		db           *gorm.DB = settings.DB
		user_id      uuid.UUID
		friend_lists []servicemodel.Friend
		errorStatus  int
		err          error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 첫 서치할때 만 작동해야 하는 함수들
	if *limit_count < 1 {
		log.Printf("%v 유저 친구 정보 찾기 가동", user_computer_number.Computer_number)
		// 유저의 아이디와 친구목록을 가져오는 로직
		if errorStatus, err = WebsocketUserStatusGetUserIdAndFriendListFunc(c, db, user_computer_number, &user_id, &friend_lists, limit_count); err != nil {
			return errorStatus, err
		}

		// 친구가 없다면 빠져나옴
		if len(friend_lists) == 0 {
			return 0, nil
		}

		// 친구 목록을 정리 하는 로직
		var (
			wg    sync.WaitGroup
			mutex sync.Mutex
		)
		wg.Add(1)
		go WebsocketUserStatusParseFriendDataFunc(&wg, &mutex, &user_id, &friend_lists, friend_statuses)
		wg.Wait()
	}

	// status가 0개라면 빠져나옴
	if len(*friend_statuses) == 0 {
		return 0, nil
	}

	// 접속중인지 확인하는 함수들
	var (
		wg        sync.WaitGroup
		mutex     sync.Mutex
		errorList []error
	)
	wg.Add(1)
	go WebsocketUserStatusVerifyFriendConnectFunc(c, db, &wg, &mutex, friend_statuses, &errorList)
	wg.Wait()

	if len(errorList) != 0 {
		for _, errorResult := range errorList {
			log.Println("시스템 오류: ", errorResult.Error())
		}
		return http.StatusInternalServerError, errors.New("유저의 상태를 업로드 하는데 오류가 발생했습니다")
	}

	return 0, nil
}

type WebsocketUserStatus interface {
	WebsocketUserStatusGetUserIdAndFriendListFunc(c context.Context, db *gorm.DB, user_computer_number *dtos.WebsocketUserComputerNumberDto, user_id *uuid.UUID, friend_lists []servicemodel.Friend, limit_count *int) (int, error)
	WebsocketUserStatusParseFriendDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_lists *[]servicemodel.Friend, friend_statuses *[]dtos.WebsocketFriendStatusDto)
	WebsocketUserStatusVerifyFriendConnectFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, mutex *sync.Mutex, friend_statuses *[]dtos.WebsocketFriendStatusDto, errorList *[]error)
}

// 유저의 아이디와 친구목록을 가져오는 로직
func WebsocketUserStatusGetUserIdAndFriendListFunc(c context.Context, db *gorm.DB, user_computer_number *dtos.WebsocketUserComputerNumberDto, user_id *uuid.UUID, friend_lists *[]servicemodel.Friend, limit_count *int) (int, error) {

	// 유저의 아이디를 가져오는 로직
	var (
		user servicemodel.User
	)
	result := db.WithContext(c).Where("computer_number = ?", user_computer_number.Computer_number).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("시스템 오류: computer number에 해당하는 유저 테이블 찾지 못함")
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 유저의 컴퓨터 넘버를 다시 확인해주세요")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저 테이블을 찾는데 오류가 발생했습니다")
		}
	}

	// 유저의 아이디를 배정
	*user_id = user.User_id

	// 친구 리스트를 가져옴
	if result = db.WithContext(c).Where("friend_1 = ? OR friend_2 = ?", *user_id, *user_id).Find(friend_lists); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저 친구의 정보를 찾는데 오류가 발생했습니다")
		}
	}

	// 친구 목록 찾기는 이정도로만 하기
	*limit_count += 1

	return 0, nil
}

// 친구 아이디만 가져오는 로직
func WebsocketUserStatusParseFriendDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_lists *[]servicemodel.Friend, friend_statuses *[]dtos.WebsocketFriendStatusDto) {

	defer wg.Done()
	for _, friend := range *friend_lists {

		mutex.Lock()
		var friend_status dtos.WebsocketFriendStatusDto
		if *user_id == friend.Friend_1 {
			friend_status.Friend_id = friend.Friend_2
			friend_status.No_check_message_number = friend.Not_Check_message_number_1
		} else if *user_id == friend.Friend_2 {
			friend_status.Friend_id = friend.Friend_1
			friend_status.No_check_message_number = friend.Not_Check_message_number_2
		}

		*friend_statuses = append(*friend_statuses, friend_status)
		mutex.Unlock()
	}
}

// 친구가 접속중인지 확인하는 로직
func WebsocketUserStatusVerifyFriendConnectFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, mutex *sync.Mutex, friend_statuses *[]dtos.WebsocketFriendStatusDto, errorList *[]error) {

	defer wg.Done()
	for idx, friend_status := range *friend_statuses {

		mutex.Lock()

		var (
			friend servicemodel.User
		)

		// 정보 찾기
		result := db.WithContext(c).Where("user_id = ?", friend_status.Friend_id).First(&friend)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				(*friend_statuses)[idx].Status = 4
				continue
			} else {
				*errorList = append(*errorList, result.Error)
			}
		}

		// 데이터 확인
		var now = time.Now()
		if friend.Computer_number == nil {
			(*friend_statuses)[idx].Status = 0
		} else {
			if !friend.UpdatedAt.Add(5 * time.Minute).Before(now) {
				(*friend_statuses)[idx].Status = 1
			} else {
				if !friend.UpdatedAt.Add(1 * time.Hour).Before(now) {
					(*friend_statuses)[idx].Status = 2
				} else {
					(*friend_statuses)[idx].Status = 3
				}
			}
		}

		mutex.Unlock()

	}

}
