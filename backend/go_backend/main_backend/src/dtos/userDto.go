package dtos

import (
	"errors"
	"unicode"
)

// 회원가입과 연관된 dto
type UserSignUpDto struct {
	Email        string `json:"email" validate:"required,email,max=80"`
	Password     string `json:"password" validate:"required,min=6,max=16"`
	Nickname     string `json:"nickname" validate:"required,min=3,max=12"`
	Term_agree_3 bool   `json:"term_agree_3" validate:"boolean"`
}

// 회원가입 dto에서 사용되는 함수
type UserSignUp interface {
	CheckPasswordType() error
}

// 회원가입에서 사용하는 함수
func (u UserSignUpDto) CheckPasswordType() error {
	var (
		isCheckNumber      bool = false
		isCheckWord        bool = false
		isCheckSpecialWord bool = false
		isCheckPassword    bool = false
	)

	passwords := u.Password
	for _, password := range passwords {
		if unicode.IsDigit(password) {
			isCheckNumber = true
		} else if unicode.IsLetter(password) {
			isCheckWord = true
		} else if unicode.IsPunct(password) || unicode.IsSymbol(password) {
			isCheckSpecialWord = true
		}

		if isCheckNumber && isCheckWord && isCheckSpecialWord {
			isCheckPassword = true
			break
		}
	}

	if !isCheckPassword {
		return errors.New("비밀번호 형식을 지켜주시길 바랍니다")
	} else {
		return nil
	}
}

// 로그인과 관련된 dto
type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}
