package models

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken    string
	RefreshToken   string
	AccessUuid     string
	RefreshUuid    string
	AccessExpires  int64
	RefreshExpires int64
}

type AccessDetails struct {
	AccessUuid string
	UserId   uint64
}