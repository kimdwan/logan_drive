package pilemodel

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteFriendChat struct {
	gorm.Model
	Friend_id         uuid.UUID `gorm:"type:uuid;not null;"`
	Friend_chat_id    uuid.UUID `gorm:"type:uuid;not null;"`
	Send_people_id    uuid.UUID `gorm:"type:uuid;not null;"`
	Address_people_id uuid.UUID `gorm:"type:uuid;not null;"`
	Message           string    `gorm:"type:varchar(2000);not null;"`
}

// 생성시 주의 사항
func (df *DeleteFriendChat) BeforeCreate(tx *gorm.DB) error {

	// 문자 갯수 확인
	var (
		message_number int = utf8.RuneCountInString(df.Message)
	)
	if message_number < 1 || message_number > 500 {
		return errors.New("문자 메세지의 갯수가 1또는 500을 넘어섭니다")
	}

	return nil
}

// 테이블 이름
func (DeleteFriendChat) TableName() string {
	return "DELETE_LOGAN_FRIEND_CHAT_TB"
}
