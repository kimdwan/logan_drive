package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/dtos"
	"gorm.io/gorm"
)

// jwt 토큰을 검증하고 payload를 제공해주는 함수
func CheckJwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			payload         dtos.Payload
			jwt_secret_keys = []string{
				os.Getenv("JWT_ACCESS_SECRET_KEY"),
				os.Getenv("JWT_REFRESH_SECRET_KEY"),
			}
			counts      int = 0
			errorStatus int
		)

		// 쿠키에서 Authorization확인
		access_token, err := ctx.Cookie("Authorization")
		if err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "클라이언트에서 보낸 토큰이 존재하지 않습니다.",
			})
			return
		}

		// 첫번째 jwt access token 확인
		if errorStatus, err = CheckJwtConfirmFunc(access_token, jwt_secret_keys[0], &payload, &counts); counts == 0 && err != nil {
			ctx.AbortWithStatusJSON(errorStatus, gin.H{
				"error": err.Error(),
			})
			return
		}

		// refresh token 확인
		if counts == 1 {
			fmt.Println(err.Error())

			// computer number가 존재하는지 확인
			computer_number := ctx.GetHeader("User-Computer-Number")
			if computer_number == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "클라이언트에서 copmuter number를 보내지 않았습니다",
				})
				return
			}

			// 본격적으로 로직을 확인
			var (
				db            *gorm.DB = settings.DB
				user          servicemodel.User
				refresh_token string
			)
			c, cancel := context.WithTimeout(context.Background(), time.Second*100)
			defer cancel()

			// refresh token 가져오기
			if errorStatus, err = CheckJwtGetRefreshTokenFunc(c, db, computer_number, &user, &refresh_token); err != nil {
				ctx.AbortWithStatusJSON(errorStatus, gin.H{
					"error": err.Error(),
				})
				return
			}

			// refresh token 검증
			if errorStatus, err = CheckJwtConfirmFunc(refresh_token, jwt_secret_keys[1], &payload, &counts); err != nil {
				if counts > 1 {
					ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
						"error": "활동을 안한지 일주일이 넘었습니다 다시 로그인 해주세요",
					})
					return
				} else {
					ctx.AbortWithStatusJSON(errorStatus, gin.H{
						"error": err.Error(),
					})
					return
				}
			}

			// 새로운 access token을 제공하고 database에도 업데이트
			if errorStatus, err = CheckJwtMakeNewAccessTokenUpdateDatabaseFunc(c, ctx, db, &user, &payload); err != nil {
				ctx.AbortWithStatusJSON(errorStatus, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		// payload를 다음 백엔드에 전달하는 함수
		if err = CheckJwtPassPayloadToBackendFunc(ctx, &payload); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

	}
}

// 함수들
type CheckJwt interface {
	CheckJwtConfirmFunc(jwt_token string, jwt_key string, payload *dtos.Payload, counts *int) (int, error)
	CheckJwtGetRefreshTokenFunc(c context.Context, db *gorm.DB, computer_number string, user *servicemodel.User, refresh_token *string) (int, error)
	CheckJwtMakeNewAccessTokenUpdateDatabaseFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User, payload *dtos.Payload) (int, error)
	CheckJwtPassPayloadToBackendFunc(ctx *gin.Context, payload *dtos.Payload) error
}

// jwt 토큰을 검증하는 함수
func CheckJwtConfirmFunc(jwt_token string, jwt_key string, payload *dtos.Payload, counts *int) (int, error) {

	// refresh token도 건너갔는지 확인하는 로직
	if *counts > 1 {
		return http.StatusUnauthorized, errors.New("refresh 토큰에서도 검증하지 못했습니다")
	}

	// 본격적으로 확인하는 로직
	jwt_value_token, err := jwt.Parse(jwt_token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt 토큰의 암호화 방식을 다시 확인해주세요")
		}

		return []byte(jwt_key), nil
	})

	// 오류 확인 (refresh 검증으로 가게 하는 로직도 존재)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			*counts += 1
			return http.StatusNoContent, errors.New("access token의 검증 시간이 지났습니다 refresh 토큰 검증으로 갑니다")
		} else {
			fmt.Println("시스템 오류: ", err.Error())
			return http.StatusUnauthorized, errors.New("jwt 토큰을 검증하는데 오류가 발생했습니다")
		}
	}

	// payload 파싱
	if claims, ok := jwt_value_token.Claims.(jwt.MapClaims); ok {
		if payload_interface, exists := claims["payload"]; exists {

			payload_byte, err := json.Marshal(payload_interface)
			if err != nil {
				fmt.Println("시스템 오류: ", err.Error())
				return http.StatusInternalServerError, errors.New("payload를 json화 하는데 오류가 발생했습니다")
			}

			if err = json.Unmarshal(payload_byte, payload); err != nil {
				fmt.Println("시스템 오류: ", err.Error())
				return http.StatusUnauthorized, errors.New("payload를 역직렬화 하는데 오류가 발생했습니다")
			}

		} else {
			return http.StatusUnauthorized, errors.New("payload로 변환하는데 오류가 발생했습니다")
		}
	} else {
		return http.StatusUnauthorized, errors.New("payload안에 claim값이 존재하지 않습니다")
	}

	return 0, nil
}

// refresh토큰을 발급받는 함수
func CheckJwtGetRefreshTokenFunc(c context.Context, db *gorm.DB, computer_number string, user *servicemodel.User, refresh_token *string) (int, error) {

	// 해당 computer_number를 가진 유저를 찾아야 함
	if result := db.WithContext(c).Where("computer_number = ?", computer_number).First(user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusUnauthorized, errors.New("클라이언트에서 보낸 computer number는 존재하지 않습니다")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 computer number에 해당하는 유저를 찾는데 오류가 발생했습니다")
		}
	}

	*refresh_token = *user.Refresh_token
	if *refresh_token == "" {
		return http.StatusUnauthorized, errors.New("refresh 토큰이 존재하지 않습니다")
	}

	return 0, nil
}

// 새로운 access token을 발급받고 업데이트 하는 함수
func CheckJwtMakeNewAccessTokenUpdateDatabaseFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *servicemodel.User, payload *dtos.Payload) (int, error) {

	var (
		jwt_tokens []string
	)

	// 새로운 jwt 토큰을 만듬
	if err := payload.MakeJwtToken(&jwt_tokens); err != nil {
		return http.StatusInternalServerError, err
	}

	// 업데이트
	user.Access_token = &jwt_tokens[0]
	user.Refresh_token = &jwt_tokens[1]
	if result := db.WithContext(c).Save(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return http.StatusInternalServerError, errors.New("데이터 베이스에 유저 정보를 업로드 하는데 오류가 발생했습니다")
	}

	// 쿠키로 보냄
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", jwt_tokens[0], 24*60*60, "", "", false, true)

	return 0, nil
}

// payload를 직렬화 하고 다음 backend에 전달하는 함수
func CheckJwtPassPayloadToBackendFunc(ctx *gin.Context, payload *dtos.Payload) error {

	// payload를 직렬화
	if payload_byte, err := json.Marshal(payload); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("payload를 직렬화 하는데 오류가 발생했습니다")
	} else {
		ctx.Set("payload_byte", string(payload_byte))
		return nil
	}
}
