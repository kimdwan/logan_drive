package dtos

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 인증정보 payload와 관련된 dtos
type Sub struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type Payload struct {
	User_id uuid.UUID `json:"user_id"`
	Sub     Sub       `json:"sub"`
}

// 인증정보와 관련된 interface
type PayloadInterface interface {
	MakeJwtToken(jwt_tokens *[]string) error
}

// jwt 토큰을 만드는 로직
func (p Payload) MakeJwtToken(jwt_tokens *[]string) error {

	var (
		jwt_secret_keys = []string{
			os.Getenv("JWT_ACCESS_SECRET_KEY"),
			os.Getenv("JWT_REFRESH_SECRET_KEY"),
		}
		jwt_time_strs = []string{
			os.Getenv("JWT_ACCESS_TIME"),
			os.Getenv("JWT_REFREH_TIME"),
		}
		jwt_times []int
	)

	// 시간 파싱
	for _, jwt_time_str := range jwt_time_strs {

		if jwt_time, err := strconv.Atoi(jwt_time_str); err != nil {
			if jwt_time_str == "" {
				fmt.Println("시스템 오류: 환경변수에 jwt time를 설정하지 않았습니다")
			} else {
				fmt.Println("시스템 오류: ", err.Error())
			}
			return errors.New("jwt을 파싱하는데 오류가 발생했습니다")
		} else {
			jwt_times = append(jwt_times, jwt_time)
		}
	}

	// secret key 생성
	for idx, jwt_secret_key := range jwt_secret_keys {
		if jwt_secret_key == "" {
			fmt.Println("시스템 오류: 환경변수에 jwt token 미기입")
			return errors.New("jwt 토큰을 생성하는데 오류가 발생했습니다")
		}

		jwt_time := jwt_times[idx]

		jwt_token_str := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"payload": p,
			"exp":     time.Now().Add(time.Duration(jwt_time) * time.Second).Unix(),
		})

		if jwt_token, err := jwt_token_str.SignedString([]byte(jwt_secret_key)); err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			return errors.New("jwt 토큰을 생성하는데 오류가 발생했습니다")
		} else {
			*jwt_tokens = append(*jwt_tokens, jwt_token)
		}
	}

	return nil
}

// 이미지와 관련된 struct
type ImgDataDto struct {
	ImgBase64 string
	ImgType   string
}

type ImgData interface {
	CheckImgType() error
}

// 이미지 타입을 확인하는 함수
func (i ImgDataDto) CheckImgType() error {

	var (
		img_type_systems = strings.Split(os.Getenv("DATABASE_USER_IMG_TYPE"), ",")
		isAllowed        = false
	)

	for _, img_type := range img_type_systems {
		if i.ImgType == img_type {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.New("이미지 타입을 다시 확인해주세요")
	}

	return nil
}
