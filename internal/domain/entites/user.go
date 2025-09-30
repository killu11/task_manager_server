package entites

import (
	"task_manager_server/pkg/security"
)

type User struct {
	name         string
	hashPassword []byte
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Hash() []byte {
	return u.hashPassword
}

func NewUser(name string, hp []byte) *User {
	return &User{
		name:         name,
		hashPassword: hp,
	}
}

func (u *User) VerifyPassword(password string) bool {
	return security.VerifyPassword(u.hashPassword, password)
}
