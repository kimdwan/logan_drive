package servicemodel

import (
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrepareFriend struct {
	gorm.Model
	Prepare_id uuid.UUID `gorm:"type:uuid;unique;not null;"`
	Request_id uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_request_respond;not null;"`
	Approve_id uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_request_respond;not null;"`
	// 상태가 제대로 작성되었는지 확인
	Status string `gorm:"type:varchar(255);not null;"`
}

// 생성전 확인
func (p *PrepareFriend) BeforeCreate(tx *gorm.DB) error {

	// 기본키 생성
	if p.Prepare_id == uuid.Nil {
		p.Prepare_id = uuid.New()
	}

	// 상태 확인
	var (
		prepare_friend_statuses []string = strings.Split(os.Getenv("DATABASE_PREVIOUS_FRIEND_STATUS_TYPE"), ",")
		isStatusTypeAllowed     bool     = false
	)
	for _, friend_status := range prepare_friend_statuses {
		if p.Status == friend_status {
			isStatusTypeAllowed = true
			break
		}
	}
	if !isStatusTypeAllowed {
		var (
			errorMsg = "허용 가능한 친구 상태 타입은: "
		)
		for idx, friend_status := range prepare_friend_statuses {
			errorMsg += friend_status
			if idx != len(prepare_friend_statuses) {
				errorMsg += ", "
			}
		}
		errorMsg += " 입니다"

		return errors.New(errorMsg)
	}

	return nil
}

// 업데이트 확인
func (p *PrepareFriend) BeforeSave(tx *gorm.DB) error {

	// 상태 확인
	var (
		prepare_friend_statuses []string = strings.Split(os.Getenv("DATABASE_PREVIOUS_FRIEND_STATUS_TYPE"), ",")
		isFriendStatusAllowed   bool     = false
	)
	for _, status := range prepare_friend_statuses {
		if p.Status == status {
			isFriendStatusAllowed = true
			break
		}
	}
	if !isFriendStatusAllowed {
		var (
			errorMsg string = "허용 가능한 친구 상태 타입은: "
		)
		for idx, status := range prepare_friend_statuses {
			errorMsg += status
			if idx != len(prepare_friend_statuses) {
				errorMsg += ", "
			}
		}
		errorMsg += " 입니다"
		return errors.New(errorMsg)
	}

	return nil
}

// 테이블 이름
func (PrepareFriend) TableName() string {
	return "LOGAN_PREPARE_FRIEND_TB"
}
