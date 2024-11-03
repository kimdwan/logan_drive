package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/src/dtos"
	"github.com/kimdwan/logan_drive/src/pkgs/services"
)

type AuthController interface {
	AuthGetUserEmailAndNickNameController(ctx *gin.Context)
	AuthGetUserProfileImgController(ctx *gin.Context)
}

// 유저의 이메일과 닉네임을 가져오는 로직
func AuthGetUserEmailAndNickNameController(ctx *gin.Context) {

	var (
		payload              *dtos.Payload
		userNicknameAndEmail *dtos.AuthNicknameAndEmailDto
		err                  error
	)

	// payload를 파싱하는 함수
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 이메일과 닉네임을 가져오는 함수
	userNicknameAndEmail = services.AuthGetUserEmailAndNickNameService(payload)

	ctx.JSON(http.StatusOK, userNicknameAndEmail)
}

// 프로필을 가져오는 로직
func AuthGetUserProfileImgController(ctx *gin.Context) {
	var (
		payload            *dtos.Payload
		user_profile_datas dtos.ImgDataDto
		errorStatus        int
		err                error
	)

	// payload를 파싱하는 함수
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		return
	}

	// profile 이미지 데이터를 가져오는 함수
	if errorStatus, err = services.AuthGetUserProfileImgService(payload, &user_profile_datas); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user_profile_datas)
}
