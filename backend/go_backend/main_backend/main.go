package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/middlewares"
	"github.com/kimdwan/logan_drive/src/pkgs/routes"
)

func init() {
	settings.LoadDotenv()
	settings.LoadDatabase()
	settings.MigrateDatabase()
}

func main() {
	port := os.Getenv("GO_PORT")
	if port == "" {
		panic("환경변수에 port 번호를 입력하지 않았습니다.")
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middlewares.CorsMiddleware())

	// 라우터 연결
	routes.UserRouter(router)

	router.Run(port)
}
