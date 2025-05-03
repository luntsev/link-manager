package auth

import (
	"link-manager/internal/user"
	"link-manager/pkg/di"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo di.IUserRepository
}

func NewAuthService(repo di.IUserRepository) *AuthService {
	return &AuthService{
		UserRepo: repo,
	}
}

func (service *AuthService) Register(email, password, name string) (*user.User, error) {
	_, err := service.UserRepo.FindByEmail(email)
	if err == nil {
		log.Println("user alredy exist")
		return nil, nil
	} else if err.Error() != "record not found" {
		log.Panicln(err.Error())
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := user.NewUser(email, string(hashedPass), name)
	return service.UserRepo.Create(user)
}

func (service *AuthService) Login(email, password string) (*user.User, error) {
	user, err := service.UserRepo.FindByEmail(email)
	if err != nil {
		log.Println("user not found")
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("wrong password")
		return nil, err
	}

	return user, nil
}
