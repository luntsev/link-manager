package user

import (
	"gorm.io/gorm"
	"link-manager/pkg/token"
)

type User struct {
	gorm.Model
	Email    string `json:"url" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Name     string `json:"name" gorm:"uniqueIndex"`
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
