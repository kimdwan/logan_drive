package dtos

import "github.com/google/uuid"

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
