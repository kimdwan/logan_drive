package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/dtos"
	"gorm.io/gorm"
)

type AuthService interface {
	AuthParsePayloadService(ctx *gin.Context) (*dtos.Payload, error)
	AuthGetUserEmailAndNickNameService(payload *dtos.Payload) *dtos.AuthNicknameAndEmailDto
	AuthGetUserProfileImgService(payload *dtos.Payload, user_profile_datas *dtos.ImgDataDto) (int, error)
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

// 유저의 프로필 이미지를 가져오는 함수
func AuthGetUserProfileImgService(payload *dtos.Payload, user_profile_datas *dtos.ImgDataDto) (int, error) {

	var (
		db               *gorm.DB = settings.DB
		profile_img_path string
		err              error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// DB에서 유저 정보를 가져오는 함수
	if err = AuthGetUserProfileImgCheckUserInDbFunc(c, db, payload, &profile_img_path); err != nil {
		return http.StatusInternalServerError, err
	}

	// 이미지가 없으면 nil을 준다
	if profile_img_path == "" {
		return 0, nil
	}

	// s3에서 데이터를 찾고 가져오는 함수
	if err = AuthGetUserProfileImgFindDataAndGetImgFunc(c, payload, profile_img_path, user_profile_datas); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

type AuthGetUserProfileImg interface {
	AuthGetUserProfileImgCheckUserInDbFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, profile_img_path *string) error
	AuthGetUserProfileImgFindDataAndGetImgFunc(c context.Context, payload *dtos.Payload, profile_img_path string, user_profile_datas *dtos.ImgDataDto) error
}

// DB에서 유저의 정보를 가져오는 함수
func AuthGetUserProfileImgCheckUserInDbFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, profile_img_path *string) error {

	var (
		user servicemodel.User
	)

	if result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(&user); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에서 유저 아이디에 해당하는 정보를 찾는데 오류가 발생했습니다")
	}

	if user.User_profile_img != nil {
		*profile_img_path = *user.User_profile_img
	}

	return nil
}

// s3에서 데이터를 찾고 가져오는 함수
func AuthGetUserProfileImgFindDataAndGetImgFunc(c context.Context, payload *dtos.Payload, profile_img_path string, user_profile_datas *dtos.ImgDataDto) error {

	// 파일 위치 확인하는 함수
	var (
		user_profile_server = os.Getenv("FILE_SERVER_USER_PROFILE_IMG")
	)
	profile_img_path = filepath.Join(user_profile_server, payload.User_id.String(), profile_img_path)

	// 파일 위치에서 데이터를 가져오는 함수
	var (
		s3client *s3.Client = settings.S3Client
	)

	imgfile_datas, err := s3client.GetObject(c, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(profile_img_path),
	})
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("디바이스에서 데이터를 찾는데 오류가 발생했습니다")
	}
	defer imgfile_datas.Body.Close()

	// 데이터를 가져온 후 작성하는 함수
	imgfile_byte_datas, err := io.ReadAll(imgfile_datas.Body)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("데이터를 변환하는데 오류가 발생했습니다")
	}

	base64Data := base64.StdEncoding.EncodeToString(imgfile_byte_datas)

	// 데이터를 집어넣는 함수
	profile_img_path_name_list := strings.Split(profile_img_path, ".")
	user_profile_datas.ImgBase64 = base64Data
	user_profile_datas.ImgType = strings.ToLower(profile_img_path_name_list[len(profile_img_path_name_list)-1])

	if err = user_profile_datas.CheckImgType(); err != nil {
		return err
	}

	return nil
}
