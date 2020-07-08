package forms

import "time"

type UserRegisterForm struct {
	UserType   int        `json:"userType,string"`
	UserStatus int        `json:"userStatus,string"`
	Gender     int        `json:"gender,string"`
	UserNmae   string     `json:"userName"`
	Password   string     `json:"password"`
	Tel        string     `json:"tel"`
	Email      string     `json:"email"`
	BirthDay   *time.Time `json:"birthDay"`
	Addr       string     `json:"addr"`
	Remark     string     `json:"remark"`
}

type DetailUpdateForm struct {
	Gender   int        `json:"gender,string"`
	Tel      string     `json:"tel"`
	Email    string     `json:"email"`
	BirthDay *time.Time `json:"birthDay"`
	Addr     string     `json:"addr"`
	Remark   string     `json:"remark"`
}

type UserStatusUpdateForm struct {
	UserStatus int `json:"userStatus,string"`
}

type UserTypeUpdateForm struct {
	UserType int `json:"userType,string"`
}

type PasswordUpdateForm struct {
	Password string `json:"password"`
}

type LoginForm struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}


type AuthUserForm struct {
	UserId     int    `json:"userId,string"`
}
