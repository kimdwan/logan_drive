# 백엔드 폴더 입니다.

#  폴더 구성

1. go_backend: go를 이용한 backend 폴더가 모여있음


# 배운점 

1. go 

- go에서는 타입으로 변형이 가능하고 이에대한 확인으로 ok로 확인할 수 있다.
- interface에 경우에만 변환이 가능하다
- io는 데이터를 읽고 byte로 변환이 가능하다
- sync에서 error 값을 추가할때 매개변수에 *[]error를 이용할 수 있다
- utf 문자의 글자수를 확인할떄 utf8.RuneCountInString()를 이용할 수 있다
- 배열자체를 바꾸고 싶으면 go에서는 해당 리스트 자체를 건드려야 한다
- go에 websocket에서 upgrader에 header는 서버에서 클라이언트에 보낼것을 이야기 한다.
- go에 websocket에서 dataType은 어떤 타입인지 확인할때 사용한다
- go에 websocket은 전역으로 연결하면 안되고 개별적으로 연결행 한다

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
