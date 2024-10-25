package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/dtos"
)

type AuthService interface {
	AuthParsePayloadService(ctx *gin.Context) (*dtos.Payload, error)
	AuthGetUserEmailAndNickNameService(payload *dtos.Payload) *dtos.AuthNicknameAndEmailDto
}

// payload를 제공하는 함수
func AuthParsePayloadService(ctx *gin.Context) (*dtos.Payload, error) {

	var (
		payload dtos.Payload
		err     error
	)

	payload_string := ctx.GetString("payload_byte")
	if err = json.Unmarshal([]byte(payload_string), &payload); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("payload를 역직렬화 하는데 오류가 발생했습니다")
	}

	return &payload, nil
}

// 유저의 이메일과 닉네임을 제공하는 함수
func AuthGetUserEmailAndNickNameService(payload *dtos.Payload) *dtos.AuthNicknameAndEmailDto {

	var (
		userNicknameAndEmail dtos.AuthNicknameAndEmailDto
	)
	userNicknameAndEmail.Email = payload.Sub.Email
	userNicknameAndEmail.Nickname = payload.Sub.Nickname

	return &userNicknameAndEmail
}
