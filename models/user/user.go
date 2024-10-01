package user

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserModel struct {
	Username string `json:"username" validate:"required,min=8,max=24"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type LoginUserModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type FetchUserModel struct {
	Id            uuid.UUID
	Username      string
	Email         string
	Password      string
	Followers_cnt int
	Following_cnt int
	Profile_img   *string
	Adress        *string
	Age           *string
	Country       *string
	Gender        *string
	Created_at    time.Time
	Updated_at    time.Time
}

type UpdateUserModel struct {
	Password string `json:"password" validate:"omitempty,min=8,max=24"`
	Username string `json:"username" validate:"omitempty,min=10,max=255"`
	Adress   string `json:"adress" validate:"omitempty,min=5,max=255"`
	Age      string `json:"age" validate:"omitempty"`
	Country  string `json:"country" validate:"omitempty"`
	Gender   string `json:"gender" validate:"omitempty"`
}
