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
	AuthUserLogoutController(ctx *gin.Context)
	AuthUserUploadProfileController(ctx *gin.Context)
	AuthUserGetFriendListController(ctx *gin.Context)
	AuthFriendSendMessageController(ctx *gin.Context)
	AuthFriendRequestController(ctx *gin.Context)
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

// 유저를 로그아웃 해주는 함수
func AuthUserLogoutController(ctx *gin.Context) {

	var (
		payload *dtos.Payload
		err     error
	)

	// payload를 가져옴
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 유저를 로그아웃 해줌
	if err = services.AuthUserLogoutService(ctx, payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "로그아웃 되었습니다.",
	})
}

// 유저가 이미지를 업로드 할 수 있게 해주는 함수
func AuthUserUploadProfileController(ctx *gin.Context) {

	var (
		payload     *dtos.Payload
		errorStatus int
		err         error
	)

	// payload를 가져오는 함수
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 유저의 프로필 이미지를 업로드 해주는 함수
	if errorStatus, err = services.AuthUserUploadProfileService(ctx, payload); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "이미지가 업로드 되었습니다.",
	})

}

// 유저의 친구 목록을 확인할 수 있는 함수
func AuthUserGetFriendListController(ctx *gin.Context) {

	var (
		payload      *dtos.Payload
		friend_lists []dtos.AuthUserFriendListDto
		errorStatus  int
		err          error
	)

	// payload를 가져오는 함수
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 친구 데이터를 가져오는 로직
	if errorStatus, err = services.AuthUserGetFriendListService(payload, &friend_lists); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, friend_lists)
}

// 친구에게 메세지를 보내는 로직
func AuthFriendSendMessageController(ctx *gin.Context) {

	var (
		payload     *dtos.Payload
		errorStatus int
		err         error
	)

	// payload 가져오는 로직
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 데이터 가져오기
	var (
		friend_message_dto dtos.AuthFriendSendMessageDto
	)
	if err = friend_message_dto.AuthFriendSendMessageParseAndPayloadFunc(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 메세지 보내기
	if errorStatus, err = services.AuthFriendSendMessageService(payload, &friend_message_dto); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

// 친구 요청
func AuthFriendRequestController(ctx *gin.Context) {

	var (
		payload          *dtos.Payload
		friend_email_dto *dtos.AuthFriendRequestEmailDto
		errorStatus      int
		err              error
	)

	// payload 가져오기
	if payload, err = services.AuthParsePayloadService(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// body 가져오기
	if friend_email_dto, err = services.AuthParseAndValidateBodyService[dtos.AuthFriendRequestEmailDto](ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 친구요청을 처리하는 함수
	if errorStatus, err = services.AuthFriendRequestService(payload, friend_email_dto); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "친구 요청 완료",
	})
}
