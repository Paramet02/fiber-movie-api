package services

import (
	"time"

	"github.com/paramet02/webapi/auth"
	
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Port: กำหนดฟังก์ชันที่ต้องมีใน Repository เพื่อให้ใช้ได้ในระบบ
type UserService interface {
	GetUserByEmail(Email string) (*User , error)
	GetUserByID(ID int) (*User , error)
	Register(email, password, firstName, lastName string) (*auth.TokenPairs, error)
	Login(email, password string) (*auth.TokenPairs, error)
}	