package user

import (
	"gorm.io/gorm"
	"link-manager/pkg/token"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Password string
	Name     string
}

func NewUser(email, password, name string) *User {
	user := User{
		Email:    email,
		Password: password,
		Name:     name,
	}
	if user.Password == "" {
		user.GenPassword(5)
	}
	return &user
}

func (u *User) GenPassword(n int) {
	u.Password = token.GenToken(n)
}
