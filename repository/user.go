package repository

import (
	"time"
)

type user struct {
	ID			int
	FirstName 	string
	LastName	string
	Email		string
	Password	string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUser(firstName, lastName string , email string , password string) *user {
    return &user{
        FirstName: firstName,
        LastName:  lastName,
		Email: email,
		Password: password, 
    }
}

// Port: กำหนดฟังก์ชันที่ต้องมีใน Repository เพื่อให้ใช้ได้ในระบบ
type UserRepository interface {
	GetUserByEmail(Email string) (*user , error)
	GetUserByID(ID int) (*user , error)
	Create(user *user) (*user , error)
}