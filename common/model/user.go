package model

import "time"

type User struct {
	ID       string `json:"id" gorm:":id;primaryKey"`
	Username string `json:"username" gorm:"username;type:varchar(100);not null"`
	Password string `json:"password" gorm:"column:password;type:varchar(200);not null"`
	Email    string `json:"email" gorm:"column:email;type:varchar(100);not null"`
	Level    string `json:"level" gorm:"level;type:int(4);not null"`
	APIKey   string `json:"api_key" gorm:"column:api_key;type:varchar(200);null"`
}

func (User) TableName() string {
	return "user"
}

type UserAuthRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
	Level    string `json:"level"`
	GroupID  string `json:"group_id"`
}

type AccessToken struct {
	ClientID     string    `json:"client_id"`
	UserID       string    `json:"user_id"`
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Createdtime  time.Time `json:"create_at"`
	Expiredtime  int       `json:"expire_at"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
	JWT          string    `json:"jwt"`
}
