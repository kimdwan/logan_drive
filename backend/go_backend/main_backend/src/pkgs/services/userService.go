package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/dtos"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService에 존재하는 함수 모음
type UserService interface {
	UserParseAndCheckBodyService(ctx *gin.Context) error
	UserSignUpService(userSignUpDto *dtos.UserSignUpDto) (int, error)
	UserLoginService(ctx *gin.Context, userLoginDto *dtos.UserLoginDto, computer_number *uuid.UUID, message *string) (int, error)
}

// 클라이언트에서 보낸 폼을 파싱하고 검증하는 함수
func UserParseAndCheckBodyService[T dtos.UserSignUpDto | dtos.UserLoginDto](ctx *gin.Context) (*T, error) {
	var (
		body T
		err  error
	)

	// 클라이언트에서 보낸 데이터를 읽는 함수
	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(json) 클라이언트에서 보낸 데이터를 읽는데 오류가 발생했습니다")
	}

	// 검증하는 함수
	validate := validator.New()
	if err = validate.Struct(body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(validate) 클라이언트에서 보낸 데이터를 읽는데 오류가 발생했습니다")
	}

	return &body, nil
}

// 회원 가입과 관련된 서비스
func UserSignUpService(userSignUpDto *dtos.UserSignUpDto) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 이메일, 닉네임 중복성 검증
	if errorStatus, err = UserSignUpCheckEmailAndNicknameAndPasswordFunc(c, db, userSignUpDto); err != nil {
		return errorStatus, err
	}

	// 비밀번호 해시화 후 데이터 베이스에 저장
	if err = UserSignUpCreateUserFunc(c, db, userSignUpDto); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 회원가입과 관련된 인터 페이스
type UserSignUpInterface interface {
	UserSignUpCheckEmailAndNicknameAndPasswordFunc(c context.Context, db *gorm.DB, userSignUpDto *dtos.UserSignUpDto) (int, error)
	UserSignUpCreateUserFunc(c context.Context, db *gorm.DB, userSignUpDto *dtos.UserSignUpDto) error
}

// 이메일, 닉네임 중복성 and 비밀번호 확인 검증
func UserSignUpCheckEmailAndNicknameAndPasswordFunc(c context.Context, db *gorm.DB, userSignUpDto *dtos.UserSignUpDto) (int, error) {

	var (
		check_user servicemodel.User
	)

	// 비밀번호 검증
	if err := userSignUpDto.CheckPasswordType(); err != nil {
		fmt.Println("시스템 오류: 비밀번호 형식 오류")
		return http.StatusBadRequest, err
	}

	// 이메일 중복성 검증
	if result := db.WithContext(c).Where("email = ?", userSignUpDto.Email).First(&check_user); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 이메일에 해당하는 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	} else {
		fmt.Println("시스템 오류: 이메일 중복")
		return http.StatusNotAcceptable, errors.New("이미 존재하는 이메일 입니다")
	}

	// 닉네임 중복성 검증
	if result := db.WithContext(c).Where("nickname = ?", userSignUpDto.Nickname).First(&check_user); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 닉네임에 해당하는 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	} else {
		fmt.Println("시스템 오류: 닉네임 중복")
		return http.StatusNotExtended, errors.New("이미 존재하는 닉네임 입니다")
	}

	return 0, nil
}

// 비밀번호 해시화 후 데이터 베이스에 저장
func UserSignUpCreateUserFunc(c context.Context, db *gorm.DB, userSignUpDto *dtos.UserSignUpDto) (err error) {

	// 비밀번호 해쉬화
	saltRoundsStr := os.Getenv("PASSWORD_SALT_ROUNDS")
	if saltRoundsStr == "" {
		fmt.Println("시스템 오류: 환경변수에 비밀번호 saltrounds 존재하지 않음")
	}

	saltRound, err := strconv.Atoi(saltRoundsStr)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("환경변수에 비밀번호를 숫자화 하는데 오류가 발생했습니다")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userSignUpDto.Password), saltRound)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("비밀번호를 해쉬화 하는데 오류가 발생했습니다")
	}

	// 유저 정보를 입력하고 + 유저 타이틀 부여
	var (
		new_user servicemodel.User
	)

	userTitleList := strings.Split(os.Getenv("DATABASE_USER_TITLE_SET"), ",")
	if len(userTitleList) == 0 {
		fmt.Println("시스템 오류: 환경변수에 DATABASE_USER_TITLE_SET을 작성하지 않았습니다")
	}

	new_user.Email = userSignUpDto.Email
	new_user.Hash = hash
	new_user.Nickname = userSignUpDto.Nickname
	new_user.Term_agree_3 = userSignUpDto.Term_agree_3
	new_user.User_title = userTitleList[len(userTitleList)-1]

	// 유저 정보 저장
	if result := db.WithContext(c).Create(&new_user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에 새로운 유저 정보를 저장하는데 오류가 발생했습니다")
	}

	return nil
}

// 로그인과 관련된 서비스
func UserLoginService(ctx *gin.Context, userLoginDto *dtos.UserLoginDto, computer_number *uuid.UUID, message *string) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		user        servicemodel.User
		jwt_tokens  []string
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 이메일과 비밀번호를 확인하는 로직
	if errorStatus, err = UserLoginCheckEmailAndPasswordFunc(c, db, &user, userLoginDto, message); err != nil {
		return errorStatus, err
	}

	// jwt 토큰을 만들고 computer number를 만드는 함수
	if err = UserLoginMakeJwtTokenAndComputerNumberFunc(&user, &jwt_tokens, computer_number); err != nil {
		return http.StatusInternalServerError, err
	}

	// 데이터 베이스에 jwt 토큰과 computer number를 업데이트 하고 token으로 access tokend을 보내는 함수
	if err = UserLoginUpdateDatabaseAndSendJwtTokenFunc(c, ctx, db, &user, &jwt_tokens, computer_number); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 로그인과 관련된 인터 페이스
type UserLoginInterface interface {
	UserLoginCheckEmailAndPasswordFunc(c context.Context, db *gorm.DB, user *servicemodel.User, userLoginDto *dtos.UserLoginDto, message *string) (int, error)
	UserLoginMakeJwtTokenAndComputerNumberFunc(user *servicemodel.User, jwt_tokens *[]string, computer_number *uuid.UUID) error
	UserLoginUpdateDatabaseAndSendJwtTokenFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User, jwt_tokens *[]string, computer_number *uuid.UUID) error
}

// 유저의 이메일과 비밀번호를 확인하는 함수
func UserLoginCheckEmailAndPasswordFunc(c context.Context, db *gorm.DB, user *servicemodel.User, userLoginDto *dtos.UserLoginDto, message *string) (int, error) {

	// 이메일 확인 로직
	if result := db.WithContext(c).Where("email = ?", userLoginDto.Email).First(user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: 이메일을 찾을수 없음")
			return http.StatusNotAcceptable, errors.New("이메일을 찾을수가 없습니다")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 이메일에 해당하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	// 비밀번호 확인 로직
	if err := bcrypt.CompareHashAndPassword(user.Hash, []byte(userLoginDto.Password)); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusNotExtended, errors.New("비밀번호를 틀렸습니다")
	}

	// 로그인을 했는지 확인하는 로직
	if user.Computer_number != nil {
		*message = "이미 다른 컴퓨터에 로그인 되어있습니다 기존의 컴퓨터를 로그아웃 하고 다시 로그인 합니다"
	} else {
		*message = "로그인 되었습니다"
	}

	return 0, nil
}

// JWT 토큰을 만들고 컴퓨터 NUMBER를 만들어 주는 함수
func UserLoginMakeJwtTokenAndComputerNumberFunc(user *servicemodel.User, jwt_tokens *[]string, computer_number *uuid.UUID) error {

	var (
		payload dtos.Payload
		sub     dtos.Sub
	)

	// payload를 만들고 jwt token을 만드는 함수
	sub.Email = user.Email
	sub.Nickname = user.Nickname
	payload.User_id = user.User_id
	payload.Sub = sub

	if err := payload.MakeJwtToken(jwt_tokens); err != nil {
		return err
	}

	// computer number를 만드는 함수
	*computer_number = uuid.New()

	return nil
}

// 데이터 베이스에 jwt 토큰과 computer number를 업데이트 하고 token으로 access tokend을 보내는 함수
func UserLoginUpdateDatabaseAndSendJwtTokenFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User, jwt_tokens *[]string, computer_number *uuid.UUID) error {

	// 데이터 베이스에 추가한다
	user.Access_token = &(*jwt_tokens)[0]
	user.Refresh_token = &(*jwt_tokens)[1]
	user.Computer_number = computer_number

	if result := db.WithContext(c).Save(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에 암호 정보를 업데이트 하는데 오류가 발생했습니다")
	}

	// 쿠키로 발송한다
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", string(*user.Access_token), 24*60*60, "", "", false, true)

	return nil
}
