package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	AuthUserLogoutService(ctx *gin.Context, payload *dtos.Payload) error
	AuthUserUploadProfileService(ctx *gin.Context, payload *dtos.Payload) error
	AuthUserGetFriendListService(payload *dtos.Payload, friend_lists *[]dtos.AuthUserFriendListDto) (int, error)
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

// 유저를 로그아웃 해주는 함수
func AuthUserLogoutService(ctx *gin.Context, payload *dtos.Payload) error {

	var (
		db   *gorm.DB = settings.DB
		user servicemodel.User
		err  error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 정보 검색
	if err = AuthUserLogoutFindDatabaseFunc(c, db, &user, payload); err != nil {
		return err
	}

	// 유저 정보 초기화
	if err = AuthUserLogoutRemoveTokenFunc(c, ctx, db, &user); err != nil {
		return err
	}

	return nil
}

type AuthUserLogout interface {
	AuthUserLogoutFindDatabaseFunc(c context.Context, db *gorm.DB, user *servicemodel.User, payload *dtos.Payload) error
	AuthUserLogoutRemoveTokenFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User) error
}

// 유저 정보를 검색한 후 초기화 해주는 함수
func AuthUserLogoutFindDatabaseFunc(c context.Context, db *gorm.DB, user *servicemodel.User, payload *dtos.Payload) error {

	// 유저 정보 검색
	result := db.WithContext(c).Where("user_id = ?", payload.User_id.String()).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("데이터 베이스에 해당 user id에 해당하는 데이터가 존재하지 않습니다")
		} else {
			log.Println("시스템 오류: ", result.Error.Error())
			return errors.New("데이터 베이스에서 유저의 정보를 찾는데 오류가 발생했습니다")
		}
	}

	return nil
}

// 보안 정보 초기화
func AuthUserLogoutRemoveTokenFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User) error {

	// 유저 정보에서 삭제
	user.Access_token = nil
	user.Refresh_token = nil
	user.Computer_number = nil
	if result := db.WithContext(c).Save(user); result.Error != nil {
		log.Println("시스템에 오류: ", result.Error.Error())
		return errors.New("데이터 베이스를 업데이트 하는데 오류가 발생했습니다")
	}

	// 토큰 전달
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", "", 0, "", "", false, true)

	return nil
}

// 유저의 프로필 이미지를 업로드 해주는 함수
func AuthUserUploadProfileService(ctx *gin.Context, payload *dtos.Payload) (int, error) {

	var (
		db                    *gorm.DB = settings.DB
		user                  servicemodel.User
		original_profile_path string
		err                   error
		errorStatus           int
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 정보를 찾는 로직
	if err = AuthUserUploadProfileFindUserAndGetOriginalPathFunc(c, db, payload, &user, &original_profile_path); err != nil {
		return http.StatusInternalServerError, err
	}

	// aws s3와 연락
	var (
		s3client                *s3.Client = settings.S3Client
		file_server_bucket_name            = os.Getenv("AWS_BUCKET_NAME")
		user_profile_img_drive             = os.Getenv("FILE_SERVER_USER_PROFILE_IMG")
	)

	// 기존의 유저가 업로드한 프로필 이미지가 있을때
	if original_profile_path != "" {
		if err = AuthUserUploadProfileOriginalProfileImgMoveUpDummyFunc(c, s3client, &file_server_bucket_name, &user_profile_img_drive, payload.User_id.String(), &original_profile_path, db, &user); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// 유저의 프로필 이미지를 업로드 해주는 로직
	if errorStatus, err = AuthUserUploadProfileCreateNewImgFunc(c, ctx, s3client, &file_server_bucket_name, &user_profile_img_drive, db, &user); err != nil {
		return errorStatus, err
	}

	return 0, nil
}

type AuthUserUploadProfile interface {
	AuthUserUploadProfileFindUserAndGetOriginalPathFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, user *servicemodel.User, original_profile_path *string) error
	AuthUserUploadProfileOriginalProfileImgMoveUpDummyFunc(c context.Context, s3client *s3.Client, file_server_bucket_name *string, user_profile_img_drive *string, user_id string, original_profile_path *string, db *gorm.DB, user *servicemodel.User) error
	AuAuthUserUploadProfileCreateNewImgFunc(c context.Context, ctx *gin.Context, s3client *s3.Client, file_server_bucket_name *string, user_profile_img_drive *string, db *gorm.DB, user *servicemodel.User) (int, error)
}

// 데이터 베이스에서 유저의 정보를 찾는 함수
func AuthUserUploadProfileFindUserAndGetOriginalPathFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, user *servicemodel.User, original_profile_path *string) error {

	// 유저 정보 찾기
	if result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(user); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에서 유저의 정보를 찾는데 오류가 발생했습니다")
	} else {
		if user.User_profile_img != nil {
			*original_profile_path = *user.User_profile_img
		}
	}

	return nil
}

// 기존의 유저의 프로필 이미지가 존재한다면 dummy 데이터로 이동해주는 함수
func AuthUserUploadProfileOriginalProfileImgMoveUpDummyFunc(c context.Context, s3client *s3.Client, file_server_bucket_name *string, user_profile_img_drive *string, user_id string, original_profile_path *string, db *gorm.DB, user *servicemodel.User) error {

	var (
		dummy_file_server_bucket_name = os.Getenv("AWS_DUMMY_BUCKET_NAME")
		err                           error
	)

	// bucket 이름이 설정 되었는지 확인하는 함수
	if dummy_file_server_bucket_name == "" {
		return errors.New("환경 변수에 bucket 이름을 설정하지 않았습니다")
	}

	// 기존 유저의 이미지 파일 위치 가져오기
	origin_user_profile_path := filepath.Join(*user_profile_img_drive, user_id, *original_profile_path)

	output_result, err := s3client.GetObject(c, &s3.GetObjectInput{
		Bucket: aws.String(*file_server_bucket_name),
		Key:    aws.String(origin_user_profile_path),
	})

	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("aws에서 기존의 이미지 파일을 찾는데 오류가 발생했습니다")
	}

	defer output_result.Body.Close()

	// 기존 파일 옮기기
	var (
		dummy__profile_server_name = os.Getenv("DUMMY_FILE_SERVER_USER_PROFILE_IMG")
		current_time_format        = time.Now().Format("2006-01-02 15:04:05")
	)

	user_original_img_file, err := io.ReadAll(output_result.Body)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("유저에게서 받아온 프로필 이미지를 byte화 하는데 오류가 발생했습니다")
	}

	user_profile_dummy_path := filepath.Join(dummy__profile_server_name, user_id, current_time_format, *original_profile_path)
	if _, err = s3client.PutObject(c, &s3.PutObjectInput{
		Bucket: aws.String(dummy_file_server_bucket_name),
		Key:    aws.String(user_profile_dummy_path),
		Body:   bytes.NewReader(user_original_img_file),
	}); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("s3에 더미 데이터를 업로드 하는데 오류가 발생했습니다")
	}

	// 기존 파일 삭제하기
	if _, err = s3client.DeleteObject(c, &s3.DeleteObjectInput{
		Bucket: aws.String(*file_server_bucket_name),
		Key:    aws.String(origin_user_profile_path),
	}); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return errors.New("s3에서 기존의 유저 프로필 이미지를 삭제하는데 오류가 발생했습니다")
	}

	// 기존 db업로드 하기
	user.User_profile_img = nil
	if result := db.WithContext(c).Save(user); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에 유저의 프로필 이름을 수정하는데 오류가 발생했습니다")
	}

	return nil
}

// 유저가 자신의 프로필 이미지를 업로드 함
func AuthUserUploadProfileCreateNewImgFunc(c context.Context, ctx *gin.Context, s3client *s3.Client, file_server_bucket_name *string, user_profile_img_drive *string, db *gorm.DB, user *servicemodel.User) (int, error) {

	// 클라이언트에서 보낸 이미지 파일 가져오기
	var (
		client_transfer_headers *multipart.FileHeader
		err                     error
	)
	if client_transfer_headers, err = ctx.FormFile("user_profile_img"); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("클라이언트에서 보낸 폼데이터에 header를 읽는 중에 오류가 발생했습니다")
	}

	// 데이터 크기와 파일의 meme type 검증
	var (
		system_approve_img_file_meme_types []string = strings.Split(os.Getenv("DATABASE_USER_IMG_TYPE"), ",")
		meme_type_allow                    bool     = false
	)

	file_name := client_transfer_headers.Filename

	// meme 타입부터 검증
	file_name_lists := strings.Split(file_name, ".")
	file_meme_type := strings.ToLower(file_name_lists[len(file_name_lists)-1])

	for _, system_approve_img_file_meme_type := range system_approve_img_file_meme_types {
		if file_meme_type == system_approve_img_file_meme_type {
			meme_type_allow = true
			break
		}
	}

	if !meme_type_allow {
		log.Println("시스템 오류: 파일의 타입이 잘못되었음")
		return http.StatusBadRequest, errors.New("파일의 타입을 다시 확인해주세요")
	}

	// 파일의 크기 검증
	var (
		file_permit_size_str = os.Getenv("FILE_PERMIT_SIZE")
		file_permit_size     int
	)
	if file_permit_size, err = strconv.Atoi(file_permit_size_str); err != nil {
		return http.StatusInternalServerError, errors.New("문자로 된 파일의 사이즈를 숫자로 바꾸는데 오류가 발생했습니다")
	}

	if client_transfer_headers.Size > int64(file_permit_size)*1024*1024 {
		log.Println("파일의 크기가 10mb보다 큼")
		return http.StatusBadRequest, errors.New("파일의 사이즈는 10mb를 넘을수가 없습니다")
	}

	// 파일을 이제 옮겨줄 예정
	client_file_data, err := client_transfer_headers.Open()
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("클라이언트에서 보낸 데이터를 읽는데 오류가 발생했습니다")
	}

	defer client_file_data.Close()
	user_profile_io_datas, err := io.ReadAll(client_file_data)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("클라이언트에서 보낸 데이터를 byte화 하는데 오류가 발생했습니다")
	}

	// 저장할 공간 만드는 중
	user_profile_img_path := filepath.Join(*user_profile_img_drive, user.User_id.String(), file_name)

	if _, err = s3client.PutObject(c, &s3.PutObjectInput{
		Bucket: aws.String(*file_server_bucket_name),
		Key:    aws.String(user_profile_img_path),
		Body:   bytes.NewReader(user_profile_io_datas),
	}); err != nil {
		log.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("aws에 s3에 데이터를 업로드 하는데 오류가 발생했습니다")
	}

	// 데이터 베이스에 유저의 정보를 업로드 해야 함
	user.User_profile_img = &file_name
	if result := db.WithContext(c).Save(user); result.Error != nil {
		log.Println("시스템 오류: ", result.Error.Error())
		return http.StatusInternalServerError, errors.New("유저의 데이터 베이스를 수정하는 데 오류가 발생했습니다")
	}

	return 0, nil
}

// 유저의 친구리스트를 확인하는 함수
func AuthUserGetFriendListService(payload *dtos.Payload, friend_lists *[]dtos.AuthUserFriendListDto) (int, error) {

	var (
		db  *gorm.DB = settings.DB
		err error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 정보를 찾는 함수
	if err = AuthUserGetFriendListFindFriendListsFunc(c, db, payload, friend_lists); err != nil {
		return http.StatusInternalServerError, err
	}

	// 친구창이 비어있으면 자동으로 나가게 설계
	if len(*friend_lists) == 0 {
		return 0, nil
	}

	// 데이터 가져오기 함수
	if err = AuthUserGetFriendListCheckUserAndGetDataFunc(c, db, friend_lists); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

type AuthUserGetFriendList interface {
	AuthUserGetFriendListFindFriendListsFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, friend_lists *[]dtos.AuthUserFriendListDto) error
	AuthUserGetFriendListFindFriendListsAddUserAsyncFunc(wg *sync.WaitGroup, mutex *sync.Mutex, payload *dtos.Payload, freind_lists *[]servicemodel.Friend, friend_list *[]dtos.AuthUserFriendListDto, errs *[]error)
	AuthUserGetFriendListCheckUserAndGetDataFunc(c context.Context, db *gorm.DB, friend_lists *[]dtos.AuthUserFriendListDto) error
	AuthUserGetFriendListCheckUserAndGetDataSyncFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, s3client *s3.Client, friend_lists *[]dtos.AuthUserFriendListDto, errs *[]error)
}

// 유저의 친구 정보가 있는지 확인하는 함수
func AuthUserGetFriendListFindFriendListsFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, friend_lists *[]dtos.AuthUserFriendListDto) error {
	var (
		friend_list []servicemodel.Friend
	)

	// 친구리스트 확인해보기
	result := db.WithContext(c).Where("friend_1 = ? OR friend_2 = ?", payload.User_id, payload.User_id).Find(&friend_list)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("시스템 오류: ", result.Error.Error())
			return errors.New("데이터 베이스에서 친구와 관련된 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	// 친구 리스트가 비어 있으면 자동으로 나오도록 함
	if len(friend_list) == 0 {
		return nil
	}

	// 친구 리스트에 따른 값을 찾는데 사용
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		errs  []error
	)
	wg.Add(1)
	go AuthUserGetFriendListFindFriendListsAddUserAsyncFunc(&wg, &mutex, payload, &friend_list, friend_lists, &errs)
	wg.Wait()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("시스템 오류: ", err.Error())
		}
		return errors.New("친구 정보를 리스트로 담는데 오류가 발생했습니다")
	}

	return nil
}

// 친구 리스트에 값을 비동기로 찾게 해주는 함수
func AuthUserGetFriendListFindFriendListsAddUserAsyncFunc(wg *sync.WaitGroup, mutex *sync.Mutex, payload *dtos.Payload, freind_lists *[]servicemodel.Friend, friend_list *[]dtos.AuthUserFriendListDto, errs *[]error) {
	defer wg.Done()

	for idx, friend_value := range *freind_lists {
		var (
			friend dtos.AuthUserFriendListDto
		)
		mutex.Lock()
		if friend_value.Friend_1 != payload.User_id && friend_value.Friend_2 == payload.User_id {
			friend.Friend_id = friend_value.Friend_1
			friend.Friend_like = friend_value.Friend_2_like
		} else {
			if friend_value.Friend_2 != payload.User_id && friend_value.Friend_1 == payload.User_id {
				friend.Friend_id = friend_value.Friend_2
				friend.Friend_like = friend_value.Friend_1_like
			} else {
				*errs = append(*errs, fmt.Errorf("%v에 index에 추가가 불가능한 부분이 존재합니다", idx))
			}
		}
		*friend_list = append(*friend_list, friend)
		mutex.Unlock()
	}

}

// 유저 정보에 맞는 데이터를 찾고 데이터 가져오기
func AuthUserGetFriendListCheckUserAndGetDataFunc(c context.Context, db *gorm.DB, friend_lists *[]dtos.AuthUserFriendListDto) error {

	var (
		wg       sync.WaitGroup
		s3client *s3.Client = settings.S3Client
		errs     []error
	)

	wg.Add(1)
	go AuthUserGetFriendListCheckUserAndGetDataSyncFunc(c, db, &wg, s3client, friend_lists, &errs)
	wg.Wait()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		return errors.New("데이터 베이스에서 친구 데이터를 만드는데 오류가 발생했습니다")
	}

	return nil
}

// 비동기로 친구 데이터 처리하기
func AuthUserGetFriendListCheckUserAndGetDataSyncFunc(c context.Context, db *gorm.DB, wg *sync.WaitGroup, s3client *s3.Client, friend_lists *[]dtos.AuthUserFriendListDto, errs *[]error) {

	defer wg.Done()

	for idx := range *friend_lists {

		err := (*friend_lists)[idx].FindUserDataAndWriteFunc(c, db, s3client)

		if err != nil {
			*errs = append(*errs, err)
		}
	}

}
