package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 기본적인 url을 파싱
		var (
			origin string = ctx.GetHeader("Origin")
		)

		// url을 파싱한다
		url_host, err := url.Parse(origin)
		if err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			fmt.Println("origin url을 파싱하는데 오류가 발생했습니다.")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		url_name := url_host.Hostname()

		// 검증
		var (
			allowed_hosts []string = strings.Split(os.Getenv("GO_ALLOWED_HOST_NAME"), ",")
			isAllowed     bool     = false
		)
		for _, allowed_host := range allowed_hosts {
			if url_name == allowed_host {
				isAllowed = true
				break
			}
		}

		// 권한 부여
		if isAllowed {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, User-Computer-Number")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// OPTIONS method는 허용하지 않음
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()

	}
}
