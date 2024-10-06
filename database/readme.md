# 정보 

- 사용 엔진: postgres 16.2

# 개발 일지 

1. 1004 
- user, deleteUser 모델 생성

# 배운점

1. 유저의 모든 권한을 가져가는 명령어
REVOKE ALL PRIVILEGES ON SCHEMA "스키마이름" FROM "user이름"

2. 데이터베이스에서 정보를 수정하는 명령어
UPDATE "DATABASE"
SET table_name = ''
WHERE 조건 = '';
