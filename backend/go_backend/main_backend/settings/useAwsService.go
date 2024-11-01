package settings

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	S3Client *s3.Client
)

func UseAwsService() {

	// 리전, 퍼블릭 키와 비밀번호 존재하는지 확인
	var (
		region_name     = os.Getenv("AWS_REGION_NAME")
		public_key      = os.Getenv("AWS_PUBLIC_KEY")
		public_password = os.Getenv("AWS_PUBLIC_PASSWORD")
	)

	if region_name == "" || public_key == "" || public_password == "" {
		panic("환경변수에 리전이름 또는 퍼블릭 키 또는 비밀번호를 입력하지 않았습니다")
	}

	// 설정 불러오기
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	provider := credentials.NewStaticCredentialsProvider(public_key, public_password, "")

	cfg, err := config.LoadDefaultConfig(c,
		config.WithRegion(region_name),
		config.WithCredentialsProvider(provider),
	)
	if err != nil {
		log.Println("시스템 오류: ", err.Error())
		panic("s3의 서비스를 로드하는데 오류가 발생했습니다")
	}

	S3Client = s3.NewFromConfig(cfg)

}
