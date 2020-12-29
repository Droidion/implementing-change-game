package models

type User struct {
	Id       uint64 `json:"id"`
	Login    string `json:"login"`
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
	UserId     uint64
}
