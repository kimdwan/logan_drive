# 백엔드 폴더 입니다.

#  폴더 구성

1. go_backend: go를 이용한 backend 폴더가 모여있음


# 배운점 

1. go 

- go에서는 타입으로 변형이 가능하고 이에대한 확인으로 ok로 확인할 수 있다.
- interface에 경우에만 변환이 가능하다

2. aws 

- aws를 코드에서 사용하기 위해서는 퍼블릭키가 필요하다 (ec2, s3 등등)
- go에서 aws를 사용하기 위해서 필요한 패키지이다.
go get - u github.com/aws/aws-sdk-go-v2/aws
go get - u github.com/aws/aws-sdk-go-v2/config
go get - u github.com/aws/aws-sdk-go-v2/service/s3
- aws를 사용하기 위해서는 config를 설정해야 한다 
config.LoadDefaultConfig(
  context, 
  리전이름, 
  인증서 (credentials)
)
- s3에 값을 입력한다 (S3Client = s3.NewFromConfig(cfg))