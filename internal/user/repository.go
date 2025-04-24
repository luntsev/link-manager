package user

import (
	"link-manager/pkg/db"
)

type UserRepository struct {
	DataBase *db.Db
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.DataBase.DB.Create(user)
	return user, result.Error
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.DataBase.DB.First(&user, "email = ?", email)
	return &user, result.Error
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: database,
	}
}
