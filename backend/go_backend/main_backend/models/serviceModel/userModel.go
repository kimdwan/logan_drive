package servicemodel

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 유저 모델
type User struct {
	gorm.Model
	User_id          uuid.UUID `gorm:"type:uuid;unique;"`
	Email            string    `gorm:"type:varchar(255);not null;unique;"`
	Hash             []byte    `gorm:"type:bytea;not null;"`
	Nickname         string    `gorm:"type:varchar(255);not null;unique;"`
	Birthday         *string   `gorm:"type:date;"`
	User_profile_img *string
	Term_agree_3     bool       `gorm:"not null;"`
	User_title       string     `gorm:"type:varchar(255);not null;"`
	Access_token     *string    `gorm:"unique;"`
	Refresh_token    *string    `gorm:"unique;"`
	Computer_number  *uuid.UUID `gorm:"type:uuid;unique;"`
}

// 생성할 때
func (u *User) BeforeCreate(tx *gorm.DB) error {

	// 유저 아이디를 부여
	if u.User_id == uuid.Nil {
		u.User_id = uuid.New()
	}

	// 이메일 확인
	var (
		emailReg string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	)
	emailCheckText := regexp.MustCompile(emailReg)
	isEmailAllowed := emailCheckText.MatchString(u.Email)
	if !isEmailAllowed {
		return errors.New("시스템 오류: 이메일 형식에 어긋남")
	}

	// 닉네임 확인
	var (
		isNicknameAllowed bool = false
	)
	nicknameLength := utf8.RuneCountInString(u.Nickname)
	if nicknameLength >= 3 && nicknameLength <= 12 {
		isNicknameAllowed = true
	}

	if !isNicknameAllowed {
		return errors.New("시스템 오류: 닉네임은 최소 3글자 최대 12글자 입니다")
	}

	// 유저 타이틀 확인
	var (
		isUserTitleAllowed bool     = false
		userTitleSets      []string = strings.Split(os.Getenv("DATABASE_USER_TITLE_SET"), ",")
	)
	for _, userTitleSet := range userTitleSets {
		if u.User_title == userTitleSet {
			isUserTitleAllowed = true
			break
		}
	}
	if !isUserTitleAllowed {
		var (
			errorMsg string = "유저 타이틀은 "
		)
		for idx, userTitleSet := range userTitleSets {
			errorMsg += userTitleSet
			if idx != len(userTitleSets)-1 {
				errorMsg += ", "
			}
		}
		errorMsg += "중 하나입니다"
		return errors.New("시스템 오류: " + errorMsg)
	}

	return nil
}

// 업데이트 할때
func (u *User) BeforeSave(tx *gorm.DB) error {

	// 닉네임 확인
	var (
		isNickNameAllowed bool = false
	)
	nicknameLength := utf8.RuneCountInString(u.Nickname)
	if nicknameLength >= 3 && nicknameLength <= 12 {
		isNickNameAllowed = true
	}
	if !isNickNameAllowed {
		return errors.New("시스템 오류: 닉네임은 최소3글자 이고 최대12글자 입니다")
	}

	// 날짜 확인
	if u.Birthday != nil {
		if _, err := time.Parse("2006-01-02", *u.Birthday); err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			return errors.New("DATE 타입에 맞지 않음")
		}
	}

	// 이미지 타입 확인
	if u.User_profile_img != nil {
		var (
			profile_img_types       []string = strings.Split(os.Getenv("DATABASE_USER_IMG_TYPE"), ",")
			isProfileImgTypeAllowed bool     = false
		)
		imgNameDotLists := strings.Split(*u.User_profile_img, ".")
		imgType := strings.ToLower(imgNameDotLists[len(imgNameDotLists)-1])

		for _, profile_img_type := range profile_img_types {
			if imgType == profile_img_type {
				isProfileImgTypeAllowed = true
				break
			}
		}

		if !isProfileImgTypeAllowed {
			var (
				errorMsg string = "이미지 타입은 "
			)

			for idx, profile_img_type := range profile_img_types {
				errorMsg += profile_img_type
				if idx != len(profile_img_types)-1 {
					errorMsg += ", "
				}
			}

			errorMsg += " 중 하나입니다"

			return errors.New("시스템 오류: " + errorMsg)
		}
	}

	// 유저 타이틀 확인
	var (
		isUserTitleAllowed bool     = false
		userTitleSets      []string = strings.Split(os.Getenv("DATABASE_USER_TITLE_SET"), ",")
	)
	for _, userTitleSet := range userTitleSets {
		if u.User_title == userTitleSet {
			isUserTitleAllowed = true
			break
		}
	}

	if !isUserTitleAllowed {
		var (
			errorMsg string = "유저 타이틀은 "
		)
		for idx, userTitleSet := range userTitleSets {
			errorMsg += userTitleSet
			if idx != len(userTitleSets)-1 {
				errorMsg += ", "
			}
		}
		errorMsg += " 중 하나입니다"

		return errors.New("시스템 오류: " + errorMsg)
	}

	return nil
}

// 테이블 이름 생성
func (User) TableName() string {
	return "LOGAN_USER_TB"
}
