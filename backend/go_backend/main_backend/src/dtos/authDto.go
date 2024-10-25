package dtos

// 유저의 닉네임과 이메일을 적는 struct
type AuthNicknameAndEmailDto struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}
