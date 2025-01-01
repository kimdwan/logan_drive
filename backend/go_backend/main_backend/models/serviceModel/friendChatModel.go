package servicemodel

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendChat struct {
	gorm.Model
	Friend_chat_id        uuid.UUID `gorm:"type:uuid;unique;not null;"`
	Send_people_id        uuid.UUID `gorm:"type:uuid;not null;"`
	Address_people_id     uuid.UUID `gorm:"type:uuid;not null;"`
	Message               string    `gorm:"type:varchar(2000);not null;"`
	Text_get_people_check int       `gorm:"type:int;not null;"`
	Whether_delete        bool      `gorm:"type:boolean;default:false;not null;"`

	// 외래키 종류
	Friend_id uuid.UUID `gorm:"type:uuid;not null;"`
}

// 생성 되기 전 확인할 함수
func (f *FriendChat) BeforeCreate(tx *gorm.DB) error {

	// 기본키가 설정되어야 한다.
	if f.Friend_chat_id == uuid.Nil {
		f.Friend_chat_id = uuid.New()
	}

	// 문자 갯수 확인
	var (
		messagenNumber int = utf8.RuneCountInString(f.Message)
	)
	if messagenNumber < 1 || messagenNumber > 500 {
		return errors.New("문자의 갯수가 500개를 넘거나 아예 없습니다")
	}

	// 최소 갯수 0 이상 확인
	if f.Text_get_people_check < 0 {
		return errors.New("문자 확인 갯수는 0보다 작을 수 없습니다")
	}

	return nil
}

// 업데이트 될때 확인할 함수
func (f *FriendChat) BeforeSave(tx *gorm.DB) error {

	// 문자 갯수 확인
	var (
		message_number int = utf8.RuneCountInString(f.Message)
	)
	if message_number < 1 || message_number > 500 {
		return errors.New("문자의 갯수가 500개를 넘거나 1개 보다 작습니다")
	}

	// 최소 갯수 0 이상 확인
	if f.Text_get_people_check < 0 {
		return errors.New("문자 확인 갯수는 0보다 작을 수 없습니다")
	}

	return nil
}

// 테이블 이름
func (FriendChat) TableName() string {
	return "LOGAN_FRIEND_CHAT_TB"
}
