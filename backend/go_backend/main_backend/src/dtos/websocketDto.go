package dtos

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// 웹소켓 에서 유저의 컴퓨터 넘버를 받는 dto
type WebsocketUserComputerNumberDto struct {
	Computer_number uuid.UUID `json:"computer_number" validate:"required,uuid"`
}

// 웹소켓에서 에러 메세지와 status를 변형 해줌
type WebsocketErrorPackDto struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// 친구의 STATUS 로그아웃 0, 접속중 1, 부재중 2(5분 이상 부재), 접속종료 3(1시간 이상 부재), 알수 없음 4
type WebsocketFriendStatusDto struct {
	Friend_id               uuid.UUID `json:"friend_id"`
	Status                  int       `json:"status"`
	No_check_message_number int       `json:"no_check_message"`
}

// 클라이언트에서 보낸 데이터 확인
type WebsocketFriendCheckDto struct {
	Computer_number uuid.UUID `json:"computer_number" validate:"required,uuid"`
	Friend_id       uuid.UUID `json:"friend_id" validate:"required,uuid"`
}

type WebsocketSendFriendMessage interface {
	WebsocketSendFriendMessageParseBodyFunc(conn *websocket.Conn) error
}

// 데이터 가져오기
func (w *WebsocketFriendCheckDto) WebsocketSendFriendMessageParseBodyFunc(conn *websocket.Conn) error {

	// 데이터 파싱
	var (
		dataType    int
		w_byte      []byte
		errorSystem WebsocketErrorPackDto
		err         error
	)
	if dataType, w_byte, err = conn.ReadMessage(); err != nil {
		log.Println("시스템 오류: ", err.Error())
		errorSystem.Error = "(json) 클라이언트 폼을 파싱하는데 오류가 발생했습니다"
		errorSystem.Status = http.StatusBadRequest
		errorMsgByte, _ := json.Marshal(errorSystem)
		conn.WriteMessage(websocket.TextMessage, errorMsgByte)
		return errors.New("(json) 클라이언트 폼을 파싱하는데 오류가 발생했습니다")
	}

	// 데이터 타입이 문제일 경우
	if dataType != websocket.TextMessage {
		errorSystem.Error = "(json) 클라이언트 폼에 데이터 타입을 확인해주세요"
		errorSystem.Status = http.StatusBadRequest
		errorMsgByte, _ := json.Marshal(errorSystem)
		conn.WriteMessage(websocket.TextMessage, errorMsgByte)
		return errors.New("(json) 클라이언트 폼에 데이터 타입을 확인해주세요")
	}

	// 데이터 적용
	if err = json.Unmarshal(w_byte, w); err != nil {
		log.Println("시스템 오류: ", err.Error())
		errorSystem.Error = "(json) 클라이언트 폼에 데이터를 읽는데 오류가 발생했습니다"
		errorSystem.Status = http.StatusBadRequest
		errorMsgByte, _ := json.Marshal(errorSystem)
		conn.WriteMessage(websocket.TextMessage, errorMsgByte)
		return errors.New("(json) 클라이언트 폼에 데이터를 읽는데 오류가 발생했습니다")
	}

	// validate 파싱
	validate := validator.New()
	if err = validate.Struct(w); err != nil {
		log.Println("시스템 오류: ", err.Error())
		errorSystem.Error = "(validate) 클라이언트 폼을 검증하는데 오류가 발생했습니다"
		errorSystem.Status = http.StatusBadRequest
		errorMsgByte, _ := json.Marshal(errorSystem)
		conn.WriteMessage(websocket.TextMessage, errorMsgByte)
		return errors.New("(validate) 클라이언트 폼을 검증하는데 오류가 발생했습니다")
	}

	return nil
}

// 문자 메세지 틀
type WebsocketFriendMessageDto struct {
	ReadType       string    `json:"readtype"`
	Message        string    `json:"message"`
	Date           time.Time `json:"time"`
	Message_number int       `json:"message_number"`
}

// 친구 요청을 보낸 유저의 정보를 가져오는 틀
type WebsocketCheckPrepareDto struct {
	Request_id  uuid.UUID `json:"request_id"`
	Postpone_id uuid.UUID `json:"postpone_id"`
}

// 유저가 보낸 데이터를 확인하는 틀
type WebsocketStreamFriendAllowStatusDto struct {
	Friend_imgbase64 string    `json:"friend_imgbase64"`
	Friend_imgtype   string    `json:"friend_imgtype"`
	Friend_email     string    `json:"friend_email"`
	Friend_nickname  string    `json:"friend_nickname"`
	Prepare_id       uuid.UUID `json:"prepare_id"`
	Friend_title     string    `json:"friend_title"`
}
