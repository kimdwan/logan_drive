package dtos

import "github.com/google/uuid"

type WebsocketComputerNumberDto struct {
	Computer_number uuid.UUID `json:"computer_number"`
}

type WebsocketFriendDto struct {
	Friend_id     uuid.UUID `json:"friend_id"`
	Friend_status int       `json:"friend_status"`
}
