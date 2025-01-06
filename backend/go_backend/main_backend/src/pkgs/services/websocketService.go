package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	WebsocketSendErrorMsgService(conn *websocket.Conn, errorStatus int, err error)
	WebsocketUserStatusService(user_computer_number *dtos.WebsocketUserComputerNumberDto, friend_statuses *[]dtos.WebsocketFriendStatusDto, limit_count int) (int, error)
	WebsocketFriendCheckMessagesService(check_friend *dtos.WebsocketFriendCheckDto, message_datas *[]dtos.WebsocketFriendMessageDto) (int, error)
	WebsocketFriendAdmitFriendAppealService(computer_number_dto *dtos.WebsocketUserComputerNumberDto, user_datas *[]dtos.WebsocketStreamFriendAllowStatusDto) (int, error)
	WebsocketFriendConfirmPrivateService(computerNumberAndFriendIdDto *dtos.WebsocketComputerNumberAndFriendIdDto) (int, error)
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

// 에러 메세지를 반환해주는 함수
func WebsocketSendErrorMsgService(conn *websocket.Conn, errorStatus int, err error) {

	var (
		errorBox = dtos.WebsocketErrorPackDto{
			Error:  err.Error(),
			Status: errorStatus,
		}
	)

	errorBoxByte, err := json.Marshal(&errorBox)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
	}

	if err = conn.WriteMessage(websocket.TextMessage, errorBoxByte); err != nil {
		log.Println("시스템 오류: ", err.Error())
	}
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

// 문자 메세지 읽기
func WebsocketFriendCheckMessagesService(check_friend *dtos.WebsocketFriendCheckDto, message_datas *[]dtos.WebsocketFriendMessageDto) (int, error) {

	var (
		db           *gorm.DB = settings.DB
		user_id      uuid.UUID
		friend_chats []servicemodel.FriendChat
		errorStatus  int
		err          error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 정보와 친구 정보 가져오기
	if errorStatus, err = WebsocketFriendCheckMessagesFindUserAndFriendFunc(c, db, check_friend, &user_id, &friend_chats); err != nil {
		return errorStatus, err
	}

	if len(friend_chats) == 0 {
		return 0, nil
	}

	// 데이터 가져오고 업데이트 하기
	if err = WebsocketFriendCheckMessagesWantDataFindUpFunc(c, db, &user_id, &friend_chats, message_datas); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

type WebsocketFriendCheckMessages interface {
	WebsocketFriendCheckMessagesFindUserAndFriendFunc(c context.Context, db *gorm.DB, check_friend *dtos.WebsocketFriendCheckDto, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat) (int, error)
	WebsocketFriendCheckMessagesWantDataFindUpFunc(c context.Context, db *gorm.DB, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat, message_datas *[]dtos.WebsocketFriendMessageDto) error
}

// 유저 정보와 친구 정보 가져오기
func WebsocketFriendCheckMessagesFindUserAndFriendFunc(c context.Context, db *gorm.DB, check_friend *dtos.WebsocketFriendCheckDto, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat) (int, error) {

	// 유저 정보 찾기
	var (
		user servicemodel.User
	)
	result := db.WithContext(c).Where("computer_number = ?", check_friend.Computer_number).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("찾을수 없는 컴퓨터 넘버")
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 컴퓨터 넘버를 다시 확인하세요")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	}

	*user_id = user.User_id

	// 데이터 찾기
	var now = time.Now().AddDate(0, -1, 0)
	if result = db.WithContext(c).Where("friend_id = ? AND updated_at > ? AND whether_delete = ?", check_friend.Friend_id, now, false).Order("created_at ASC").Find(friend_chats); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 친구와 데화한 내용을 찾는데 오류가 발생했습니다")
		}
	}

	return 0, nil
}

// 필요한 메세지 데이터만 가져오고 저장하기
func WebsocketFriendCheckMessagesWantDataFindUpFunc(c context.Context, db *gorm.DB, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat, message_datas *[]dtos.WebsocketFriendMessageDto) error {

	// 데이터 가져오기
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	wg.Add(1)
	go WebsocketFriendCheckMessageWantDataFindUpParseDataFunc(&wg, &mutex, user_id, friend_chats, message_datas)
	wg.Wait()

	// 데이터 업데이트
	if result := db.WithContext(c).Save(friend_chats); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터를 업데이트 하는데 오류가 발생했습니다")
	}

	return nil
}

type WebsocketFriendCheckMessageWantDataFindUp interface {
	WebsocketFriendCheckMessageWantDataFindUpParseDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat, message_datas *[]dtos.WebsocketFriendMessageDto)
}

// 메세지 읽어오는 로직
func WebsocketFriendCheckMessageWantDataFindUpParseDataFunc(wg *sync.WaitGroup, mutex *sync.Mutex, user_id *uuid.UUID, friend_chats *[]servicemodel.FriendChat, message_datas *[]dtos.WebsocketFriendMessageDto) {

	defer wg.Done()
	for idx, frined_chat := range *friend_chats {
		var (
			message_data dtos.WebsocketFriendMessageDto
		)

		mutex.Lock()

		// 수신자 체크
		if *user_id == frined_chat.Send_people_id {
			message_data.ReadType = "send"
			message_data.Message_number = frined_chat.Text_get_people_check
		} else {
			message_data.ReadType = "access"
			(*friend_chats)[idx].Text_get_people_check = 0
			message_data.Message_number = 0
		}

		// 데이터 정리
		message_data.Message = frined_chat.Message
		message_data.Date = frined_chat.CreatedAt

		(*message_datas) = append((*message_datas), message_data)

		mutex.Unlock()

	}

}

// 친구요청 정보를 실시간으로 확인해주는 함수
func WebsocketFriendAdmitFriendAppealService(computer_number_dto *dtos.WebsocketUserComputerNumberDto, user_datas *[]dtos.WebsocketStreamFriendAllowStatusDto) (int, error) {

	// 유저 데이터를 초기화
	*user_datas = nil

	var (
		db          *gorm.DB = settings.DB
		CheckDatas  []dtos.WebsocketCheckPrepareDto
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 요청온 친구의 정보를 실시간으로 확인하고 체크하는 함수
	if errorStatus, err = WebsocketFriendAdmitFriendAppealDataCheckFunc(c, db, &CheckDatas, computer_number_dto); err != nil {
		return errorStatus, err
	}

	if len(CheckDatas) == 0 {
		return 0, nil
	}

	// 요청온 데이터를 확인하고 정보를 가져오는 함수
	if errorStatus, err = WebsocketFriendAdmitGetPostPoneUserProfileFunc(c, db, &CheckDatas, user_datas); err != nil {
		return errorStatus, err
	}

	return 0, nil
}

type WebsocketFriendAdmitFriendAppeal interface {
	WebsocketFriendAdmitFriendAppealDataCheckFunc(c context.Context, db *gorm.DB, CheckDatas *[]dtos.WebsocketCheckPrepareDto, computer_number_dto *dtos.WebsocketUserComputerNumberDto) (int, error)
	WebsocketFriendAdmitGetPostPoneUserProfileFunc(c context.Context, db *gorm.DB, CheckDatas *[]dtos.WebsocketCheckPrepareDto, websocketStreamDtos *[]dtos.WebsocketStreamFriendAllowStatusDto) (int, error)
}

// 요청온 친구의 정보를 실시간으로 확인하고 체크하는 함수
func WebsocketFriendAdmitFriendAppealDataCheckFunc(c context.Context, db *gorm.DB, CheckDatas *[]dtos.WebsocketCheckPrepareDto, computer_number_dto *dtos.WebsocketUserComputerNumberDto) (int, error) {

	// 유저 정보부터 확인
	var (
		user servicemodel.User
	)
	result := db.WithContext(c).Where("computer_number = ?", computer_number_dto.Computer_number).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("해당 유저를 찾을 수 없음")
			return http.StatusBadRequest, errors.New("컴퓨터 넘버가 잘못되었습니다 다시 확인해주세요")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	}

	// 친구창 온거 데이터 확인
	var (
		prepare_friends []servicemodel.PrepareFriend
	)
	if result = db.WithContext(c).Where("approve_id = ?", user.User_id).Order("created_at ASC").Find(&prepare_friends); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return http.StatusInternalServerError, errors.New("데이터 베이스에서 승인 유저를 찾는데 오류가 발생했습니다")
	}

	// 데이터 정리
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	wg.Add(1)
	go WebsocketFriendAdmitFriendAppealDataCheckOrganizeFunc(&wg, &mutex, &prepare_friends, CheckDatas)
	wg.Wait()

	return 0, nil
}

type WebsocketFriendAdmitFriendAppealDataCheck interface {
	WebsocketFriendAdmitFriendAppealDataCheckOrganizeFunc(wg *sync.WaitGroup, mutex *sync.Mutex, prepare_friends *[]servicemodel.PrepareFriend, CheckDatas *[]dtos.WebsocketCheckPrepareDto)
}

// 친구 요청을 보낸 사람들의 아이디를 가져오는 함수
func WebsocketFriendAdmitFriendAppealDataCheckOrganizeFunc(wg *sync.WaitGroup, mutex *sync.Mutex, prepare_friends *[]servicemodel.PrepareFriend, CheckDatas *[]dtos.WebsocketCheckPrepareDto) {
	defer wg.Done()
	for _, prepare_friend := range *prepare_friends {
		var (
			system_statuses []string = strings.Split(os.Getenv("DATABASE_PREVIOUS_FRIEND_STATUS_TYPE"), ",")
		)
		mutex.Lock()
		if prepare_friend.Status == system_statuses[0] {
			var prepare_id = dtos.WebsocketCheckPrepareDto{
				Request_id:  prepare_friend.Request_id,
				Postpone_id: prepare_friend.Prepare_id,
			}
			*CheckDatas = append(*CheckDatas, prepare_id)
		}
		mutex.Unlock()
	}

}

// 유저의 데이터를 가져오는 함수
func WebsocketFriendAdmitGetPostPoneUserProfileFunc(c context.Context, db *gorm.DB, CheckDatas *[]dtos.WebsocketCheckPrepareDto, websocketStreamDtos *[]dtos.WebsocketStreamFriendAllowStatusDto) (int, error) {

	var (
		s3Client *s3.Client = settings.S3Client
		wg       sync.WaitGroup
		mutex    sync.Mutex
		errorBox []error
	)

	wg.Add(1)
	go WebsocketFriendAdmitGetPostPoneUserProfileWorkOnFunc(c, &wg, &mutex, db, s3Client, CheckDatas, websocketStreamDtos, &errorBox)
	wg.Wait()

	if len(errorBox) != 0 {
		return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저 정보를 찾는데 오류가 발생했습니다")
	}

	return 0, nil
}

type WebsocketFriendAdmitGetPostPoneUserProfile interface {
	WebsocketFriendAdmitGetPostPoneUserProfileWorkOnFunc(c context.Context, wg *sync.WaitGroup, mutex *sync.Mutex, db *gorm.DB, s3Client *s3.Client, CheckDatas *[]dtos.WebsocketCheckPrepareDto, websocketStreamDtos *[]dtos.WebsocketStreamFriendAllowStatusDto, errorBox *[]error)
}

// 데이터 수집 함수
func WebsocketFriendAdmitGetPostPoneUserProfileWorkOnFunc(c context.Context, wg *sync.WaitGroup, mutex *sync.Mutex, db *gorm.DB, s3Client *s3.Client, CheckDatas *[]dtos.WebsocketCheckPrepareDto, websocketStreamDtos *[]dtos.WebsocketStreamFriendAllowStatusDto, errorBox *[]error) {

	defer wg.Done()
	for _, check_data := range *CheckDatas {
		// 유저 데이터에서 정보 찾기
		var (
			user servicemodel.User
		)
		mutex.Lock()
		result := db.WithContext(c).Where("user_id = ?", check_data.Request_id).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Println("존재하지 않는 유저")
			} else {
				log.Println("시스템 오류: ", result.Error.Error())
				*errorBox = append(*errorBox, errors.New("데이터를 수집하는데 오류가 발생했습니다"))
			}
		}

		// 데이터 정리
		var (
			websocketStreamDto dtos.WebsocketStreamFriendAllowStatusDto
		)
		websocketStreamDto.Prepare_id = check_data.Postpone_id
		websocketStreamDto.Friend_email = user.Email
		websocketStreamDto.Friend_nickname = user.Nickname
		websocketStreamDto.Friend_title = user.User_title

		// 유저의 사진 가져오기
		if user.User_profile_img != nil {
			var (
				bucket_name         string = os.Getenv("AWS_BUCKET_NAME")
				file_server_profile string = os.Getenv("FILE_SERVER_USER_PROFILE_IMG")
				profile_name               = *user.User_profile_img
			)
			file_position := filepath.Join(file_server_profile, user.User_id.String(), profile_name)
			if s3_result, err := s3Client.GetObject(c, &s3.GetObjectInput{
				Bucket: aws.String(bucket_name),
				Key:    aws.String(file_position),
			}); err != nil {
				log.Println("시스템 오류: ", err.Error())
				*errorBox = append(*errorBox, err)
			} else {
				defer s3_result.Body.Close()
				user_email_data, err := io.ReadAll(s3_result.Body)
				if err != nil {
					log.Println("시스템 오류: ", err.Error())
					*errorBox = append(*errorBox, err)
				}

				basedata := base64.StdEncoding.EncodeToString(user_email_data)
				websocketStreamDto.Friend_imgbase64 = basedata
			}

			// 프로필 데이터의 타입
			profile_data_types := strings.Split(profile_name, ".")
			websocketStreamDto.Friend_imgtype = profile_data_types[len(profile_data_types)-1]
		}

		*websocketStreamDtos = append(*websocketStreamDtos, websocketStreamDto)
		mutex.Unlock()
	}

}

// 유저 한명의 정보를 실시간으로 수집하기 위해서
func WebsocketFriendConfirmPrivateService(computerNumberAndFriendIdDto *dtos.WebsocketComputerNumberAndFriendIdDto, Friend_detail *dtos.WebsocketCheckFriendDetailDto) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		friend      servicemodel.User
		friend_id   uuid.UUID
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 친구의 정보를 찾고 컴퓨터 넘버가 맞다면 가져오고 아니면 반송하는 로직
	if errorStatus, err = WebsocketFriendConfirmPrivateCheckFriendIdAndCreateDatabaseFunc(c, db, &friend, computerNumberAndFriendIdDto, &friend_id); err != nil {
		return errorStatus, err
	}

	// 친구의 데이터를 가져오는 로직
	var (
		s3client *s3.Client = settings.S3Client
	)
	if err = WebsocketFriendConfirmPrivateSetoffFriendDataFunc(c, s3client, &friend, Friend_detail, &friend_id); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

type WebsocketFriendConfirmPrivate interface {
	WebsocketFriendConfirmPrivateCheckFriendIdAndCreateDatabaseFunc(c context.Context, db *gorm.DB, friend *servicemodel.User, computerNumberAndFriendIdDto *dtos.WebsocketComputerNumberAndFriendIdDto, friend_id *uuid.UUID) (int, error)
	WebsocketFriendConfirmPrivateSetoffFriendDataFunc(c context.Context, s3client *s3.Client, friend *servicemodel.User, friend_detail *dtos.WebsocketCheckFriendDetailDto, friend_id *uuid.UUID) error
}

// 친구 정보를 찾아볼 뿐만 아니라 컴퓨터 넘버가 맞는지 친구는 맞는지를 검증하는 로직
func WebsocketFriendConfirmPrivateCheckFriendIdAndCreateDatabaseFunc(c context.Context, db *gorm.DB, friend *servicemodel.User, computerNumberAndFriendIdDto *dtos.WebsocketComputerNumberAndFriendIdDto, friend_id *uuid.UUID) (int, error) {

	// 컴퓨터 넘버가 맞는지 확인하는 로직
	var (
		user servicemodel.User
	)
	result := db.WithContext(c).Where("computer_number = ?", computerNumberAndFriendIdDto.Computer_number).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("컴퓨터 넘버를 다시 확인해야 함")
			return http.StatusBadRequest, errors.New("컴퓨터 넘버를 다시 확인해주시길 바랍니다")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 컴퓨터 넘버에 해당하는 데이터 베이스를 찾는데 오류가 발생했습니다")
		}
	}

	// 친구가 맞는지 확인하는 로직
	var (
		friend_model servicemodel.Friend
	)
	if result = db.WithContext(c).Where("friend_1 =? AND friend_2 = ?", user.User_id, computerNumberAndFriendIdDto.Friend_id).Or("friend_1 = ? AND friend_2 = ?", computerNumberAndFriendIdDto.Friend_id, user.User_id).First(&friend_model); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("둘은 친구관계가 아님")
			return http.StatusBadRequest, errors.New("둘은 친구관계가 아닙니다")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 친구에 해당하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}
	*friend_id = friend_model.Friend_id

	// 친구 데이터 가져오기
	if result = db.WithContext(c).Where("user_id = ?", computerNumberAndFriendIdDto.Friend_id).First(friend); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("유저의 정보를 찾을 수 없음")
			return http.StatusBadRequest, errors.New("해당 유저는 존재하지 않는 유저 입니다")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	}

	return 0, nil
}

// 친구의 디테일한 데이터를 가져오는 로직
func WebsocketFriendConfirmPrivateSetoffFriendDataFunc(c context.Context, s3client *s3.Client, friend *servicemodel.User, friend_detail *dtos.WebsocketCheckFriendDetailDto, friend_id *uuid.UUID) error {

	// 기본적으로 배정할 수 있는 데이터 배정
	friend_detail.Friend_email = friend.Email
	friend_detail.Friend_id = *friend_id
	friend_detail.Friend_nickname = friend.Nickname
	friend_detail.Friend_title = friend.User_title

	// 이미지 관련
	if friend.User_profile_img != nil {
		var (
			bucket_name        string = os.Getenv("AWS_BUCKET_NAME")
			img_server_postion string = os.Getenv("FILE_SERVER_USER_PROFILE_IMG")
			profile_img        string = *friend.User_profile_img
		)
		profile_img_path := filepath.Join(img_server_postion, friend.User_id.String(), profile_img)

		s3_output, err := s3client.GetObject(c, &s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(profile_img_path),
		})
		if err != nil {
			log.Println("시스템 오류: ", err.Error())
			return errors.New("s3에서 데이터를 찾는데 오류가 발생했습니다")
		}
		defer s3_output.Body.Close()

		s3_profile_data, err := io.ReadAll(s3_output.Body)
		if err != nil {
			log.Println("시스템 오류: ", err.Error())
			return errors.New("데이터를 byte화 하는데 오류가 발생했습니다")
		}

		s3_base_64_data := base64.StdEncoding.EncodeToString(s3_profile_data)

		friend_detail.Friend_imgbase64 = s3_base_64_data

		// 타입 지정
		var (
			profile_img_lists []string = strings.Split(profile_img, ".")
		)
		friend_detail.Freind_imgtype = profile_img_lists[len(profile_img_lists)-1]
	}

	// 유저의 status 검증
	if friend.Computer_number != nil {
		var (
			friend_time           = friend.UpdatedAt
			now_datas   time.Time = time.Now()
			sub_datas             = now_datas.Sub(friend_time)
		)

		if sub_datas <= 5*time.Minute {
			friend_detail.Status = 1
		} else if sub_datas <= 1*time.Hour {
			friend_detail.Status = 2
		} else {
			friend_detail.Status = 3
		}

	} else {
		friend_detail.Status = 0
	}

	return nil

}
