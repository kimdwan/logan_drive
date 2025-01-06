package dtos

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"gorm.io/gorm"
)

// 유저의 닉네임과 이메일을 적는 struct
type AuthNicknameAndEmailDto struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// 유저의 친구리스트를 보는 struct
type AuthUserFriendListDto struct {
	Friend_imgbase64 string    `json:"friend_imgbase64,omitempty"`
	Friend_imgtype   string    `json:"friend_imgtype,omitempty"`
	Friend_email     string    `json:"friend_email"`
	Friend_nickname  string    `json:"friend_nickname"`
	Friend_id        uuid.UUID `json:"friend_id"`
	Friend_like      bool      `json:"friend_like"`
	Friend_title     string    `json:"friend_title"`
}

type AuthUserFriend interface {
	FindUserDataAndWriteFunc(c context.Context, db *gorm.DB, s3client *s3.Client) error
}

// 데이터를 추가해주는 함수
func (aufl *AuthUserFriendListDto) FindUserDataAndWriteFunc(c context.Context, db *gorm.DB, s3client *s3.Client) error {

	var (
		user servicemodel.User
	)

	// 유저 데이터 찾기
	result := db.WithContext(c).Where("user_id = ?", aufl.Friend_id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("삭제 되었거나 존재하지 않는 유저")
			return nil
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return errors.New("데이터 베이스에서 friend_id에 해당하는 user 테이블을 찾는데 오류가 발생했습니다")
		}
	}

	// 기본적인 이메일과 닉네임 세팅
	aufl.Friend_email = user.Email
	aufl.Friend_nickname = user.Nickname
	aufl.Friend_title = user.User_title

	// 유저 프로필 이미지가 존재한다면
	if user.User_profile_img != nil {

		var (
			bucket_name string = os.Getenv("AWS_BUCKET_NAME")
			file_server string = os.Getenv("FILE_SERVER_USER_PROFILE_IMG")
		)

		if bucket_name == "" || file_server == "" {
			return errors.New("환경 변수에 파일 서버와 bucket 이름을 설정하지 않았습니다")
		}

		// 프로필 url 가져오기
		friend_profile_url := filepath.Join(file_server, aufl.Friend_id.String(), *user.User_profile_img)

		friend_profile_img_data, err := s3client.GetObject(c, &s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(friend_profile_url),
		})

		if err != nil {
			log.Println("시스템 오류: ", err.Error())
			return errors.New("s3에서 데이터를 찾는데 오류가 발생했습니다")
		}

		// 이미지 데이터를 파싱하고 집어넣기
		defer friend_profile_img_data.Body.Close()
		friend_profile_img, err := io.ReadAll(friend_profile_img_data.Body)
		if err != nil {
			log.Println("시스템 오류: ", err.Error())
			return errors.New("프로필 이미지를 읽는데 오류가 발생했습니다")
		}

		// base 64로 파싱하고 데이터 입력
		friend_profile_img_base64_data := base64.StdEncoding.EncodeToString(friend_profile_img)

		aufl.Friend_imgbase64 = friend_profile_img_base64_data

		// 프로필 이미지에 타입을 넣기
		var (
			img_profile_data_lists []string = strings.Split(*user.User_profile_img, ".")
		)

		aufl.Friend_imgtype = strings.ToLower(img_profile_data_lists[len(img_profile_data_lists)-1])
	}

	return nil
}

// 친구에게 메세지를 보낼때 사용하는 스트럭트
type AuthFriendSendMessageDto struct {
	Friend_id uuid.UUID `json:"friend_id" validate:"required,uuid"`
	Message   string    `json:"message" validate:"required,min=1,max=500"`
}

// 맞는지 검증하는 로직
type AuthFriendSendMessage interface {
	AuthFriendSendMessageParseAndPayloadFunc(ctx *gin.Context) error
}

func (a *AuthFriendSendMessageDto) AuthFriendSendMessageParseAndPayloadFunc(ctx *gin.Context) (err error) {

	// body에 데이터를 가져오기
	if err = ctx.ShouldBindBodyWithJSON(a); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("(json) 클라이언트에 데이터를 가져오는데 오류가 발생했습니다")
	}

	// validate 검증
	validate := validator.New()
	if err = validate.Struct(a); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("(validate) 클라이언트에 폼을 검증하는데 오류가 발생했습니다")
	}

	return nil
}

// 친구요청을 보낼때 사용하는 스트럭트
type AuthFriendRequestEmailDto struct {
	Email string `json:"email" validate:"required,email"`
}

// 친구 요청에 대해서 확인하고자 할때 사용함
type AuthFriendConfirmTypeDto struct {
	Permit_id  uuid.UUID `json:"permit_id" validate:"required,uuid"`
	Allow_type string    `json:"allow_type" validate:"required"`
}

type AuthFriendConfirmType interface {
	CheckAllowTypeFunc() error
}

// 해당 허용 타입이 환경변수에서 허용한 것인가
func (a AuthFriendConfirmTypeDto) CheckAllowTypeFunc() error {

	var (
		system_allow_types = []string{
			"allow",
			"delete",
		}
		isAllowedType bool = false
	)

	for _, allow_type := range system_allow_types {
		if a.Allow_type == allow_type {
			isAllowedType = true
			break
		}
	}

	if !isAllowedType {
		return errors.New("타입을 다시 확인해주세요")
	}

	return nil
}
