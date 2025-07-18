package domain

import "time"

type User struct {
	ID          int32
	Username    string
	Email       string
	Password    string
	LastLoginAt time.Time
}

type UserRepository interface {
	Create(user *User) (int32, error)
	GetByUsername(userName string) (*User, error)
	GetByEmail(email string) (*User, error)
}
