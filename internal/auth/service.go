package auth

import (
	"fmt"
	"link-manager/internal/user"
)

type AuthService struct {
	userRepo *user.UserRepository
}

func NewAuthService(repo *user.UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (service *AuthService) Register(email, password, name string) (*user.User, error) {
	_, err := service.userRepo.FindByEmail(email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user := user.NewUser(email, password, name)
	return service.userRepo.Create(user)
}
