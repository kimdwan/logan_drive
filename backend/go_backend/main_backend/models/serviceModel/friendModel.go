package servicemodel

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friend struct {
	gorm.Model
	Friend_id                  uuid.UUID `gorm:"type:uuid;unique;not null;"`
	Friend_1                   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_friends;not null;"`
	Friend_2                   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_friends;not null;"`
	Friend_1_like              bool      `gorm:"not null;default:false;"`
	Friend_2_like              bool      `gorm:"not null;default:false;"`
	Not_Check_message_number_1 int       `gorm:"type:int;not null;default:0;"`
	Not_Check_message_number_2 int       `gorm:"type:int;not null;default:0;"`
	// 관계 모음
	Friend_chat []FriendChat `gorm:"foreignKey:friend_id;references:friend_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// 테이블을 생성할때 확인하는 함수
func (f *Friend) BeforeCreate(tx *gorm.DB) error {

	// uuid 생성
	if f.Friend_id == uuid.Nil {
		f.Friend_id = uuid.New()
	}

	// 데이터 베이스에 실제 친구가 존재하는지 확인하는 창
	var (
		checkUser_1 User
		checkUser_2 User
	)

	if result := tx.Where("user_id = ?", f.Friend_1).First(&checkUser_1); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("데이터 베이스에 존재하지 않는 유저 아이디 입니다 (첫번째 친구)")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return errors.New("데이터 베이스에서 유저 정보를 찾는데 오류가 발생했습니다 (첫번째 친구)")
		}
	} else {
		if result = tx.Where("user_id = ?", f.Friend_2).First(&checkUser_2); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("데이터 베이스에 존재하지 않는 유저 아이디 입니다 (두번째 친구)")
			} else {
				log.Println("시스템 오류: ", result.Error.Error())
				return errors.New("데이터 베이스에서 유저 정보를 찾는데 오류가 발생했습니다 (두번째 친구)")
			}
		}
	}

	return nil
}

// 테이블을 수정할때 확인하는 함수
func (f *Friend) BeforeSave(tx *gorm.DB) error {

	if f.Not_Check_message_number_1 < 0 || f.Not_Check_message_number_2 < 0 {
		return errors.New("친구가 메세지를 확인하지 못한 메세지에 갯수는 음수가 나올수가 없습니다")
	}

	return nil
}

// 테이블 이름
func (Friend) TableName() string {
	return "LOGAN_FRIEND_TB"
}
