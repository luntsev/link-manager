package auth

import (
	"golang.org/x/crypto/bcrypt"
	"link-manager/internal/user"
	"log"
)

type AuthService struct {
	UserRepo *user.UserRepository
}

func NewAuthService(repo *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: repo,
	}
}

func (service *AuthService) Register(email, password, name string) (*user.User, error) {
	_, err := service.UserRepo.FindByEmail(email)
	if err == nil {
		log.Println("user alredy exist")
		return nil, nil
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := user.NewUser(email, string(hashedPass), name)
	return service.UserRepo.Create(user)
}
