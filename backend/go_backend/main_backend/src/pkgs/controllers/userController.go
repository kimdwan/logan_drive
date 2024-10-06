package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kimdwan/logan_drive/src/dtos"
	"github.com/kimdwan/logan_drive/src/pkgs/services"
)

// UserController에서 함수들 정리하는 장소
type UserController interface {
	UserSignUpController(ctx *gin.Context)
	UserLoginController(ctx *gin.Context)
}

// 회원가입의 컨트롤러를 담당하는 함수
func UserSignUpController(ctx *gin.Context) {

	var (
		userSignUpDto *dtos.UserSignUpDto
		errorStatus   int
		err           error
	)

	// 클라이언트에서 보낸 폼을 확인하는 함수
	if userSignUpDto, err = services.UserParseAndCheckBodyService[dtos.UserSignUpDto](ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 회원가입이 이루어지는 로직
	if errorStatus, err = services.UserSignUpService(userSignUpDto); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "회원가입 되었습니다.",
	})

}

// 로그인 컨트롤러를 담당하는 함수
func UserLoginController(ctx *gin.Context) {

	var (
		userLoginDto    *dtos.UserLoginDto
		computer_number uuid.UUID
		message         string
		errorStatus     int
		err             error
	)

	// 클라이언트에서 보낸 폼을 확인하는 함수
	if userLoginDto, err = services.UserParseAndCheckBodyService[dtos.UserLoginDto](ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 로그인이 이루어 지는 로직
	if errorStatus, err = services.UserLoginService(ctx, userLoginDto, &computer_number, &message); err != nil {
		ctx.AbortWithStatusJSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"computer_number": computer_number,
		"message":         message,
	})
}
