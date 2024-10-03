package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/logan_drive/settings"
	"github.com/kimdwan/logan_drive/src/middlewares"
)

func init() {
	settings.LoadDotenv()
	settings.LoadDatabase()
}

func main() {
	port := os.Getenv("GO_PORT")
	if port == "" {
		panic("환경변수에 port 번호를 입력하지 않았습니다.")
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middlewares.CorsMiddleware())

	router.Run(port)
}
