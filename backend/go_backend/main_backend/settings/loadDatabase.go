package settings

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	var (
		dsn string = os.Getenv("GO_DATABASE_DSN")
		err error
	)

	if dsn == "" {
		panic("환경변수에 데이터 베이스 dsn을 설정하지 않았습니다.")
	}

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		panic("데이터 베이스를 로드 하는데 오류가 발생했습니다.")
	}

}
