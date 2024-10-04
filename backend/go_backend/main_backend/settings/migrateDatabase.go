package settings

import (
	"fmt"

	pilemodel "github.com/kimdwan/logan_drive/models/pileModel"
	servicemodel "github.com/kimdwan/logan_drive/models/serviceModel"
)

func MigrateDatabase() {

	if err := DB.AutoMigrate(&servicemodel.User{}, &pilemodel.DeleteUser{}); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		panic("데이터 베이스에 마이그레이션 하는데 오류가 발생했습니다")
	}

}
