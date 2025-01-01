# 정보 

- 사용 엔진: postgres 16.2

# 개발 일지 

1. 1004 
- user, deleteUser 모델 생성

2. 1120
- friend, deleteFriend 모델 생성

3. 0101 
- friend 테이블에 유저의 메세지 읽음 갯수 추가
- friendChat, deleteFriendChat 모델 생성
- preparefriendstatus 모델 생성

# 배운점

1. 유저의 모든 권한을 가져가는 명령어
REVOKE ALL PRIVILEGES ON SCHEMA "스키마이름" FROM "user이름"

2. 데이터베이스에서 정보를 수정하는 명령어
UPDATE "DATABASE"
SET table_name = ''
WHERE 조건 = '';

3. 데이터 베이스에서 정보를 여러개 수정하기 위해서는 
UPDATE "테이블 이름" 
SET "column1" = '',
"column2" = ''
WHERE "columns" = ''

4. 데이터 베이스에서 테이블을 생성할때 기본값을 FALSE로 하려면
CREATE TABLE "table_name" (
  id INT PRIMARY KEY, 
  bool1 BOOLEAN  NOT NULL DEFAULT FALSE
)
