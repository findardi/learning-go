package model

type RequestUserRegister struct {
	Username string `json:"username" validate:"required,min=8,max=25"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type RequestUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Passowrd string `json:"password" validate:"required,min=8,max=255"`
}

type ResponseUserLogin struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ResponseUser struct {
	Id        int32  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
