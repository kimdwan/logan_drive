package pilemodel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 유저가 탈퇴했을때 저장하는 더미 데이터
type DeleteUser struct {
	gorm.Model
	User_id          uuid.UUID `gorm:"type:uuid;not null;"`
	Email            string    `gorm:"type:varchar(255);not null;"`
	Birthday         *string   `gorm:"type:date;"`
	User_profile_img *string
	Nickname         string `gorm:"type:varchar(255);not null;"`
	User_title       string `gorm:"type:varchar(255);not null;"`
}

// 테이블 이름
func (DeleteUser) TableName() string {
	return "DELETE_LOGAN_USER_TB"
}
