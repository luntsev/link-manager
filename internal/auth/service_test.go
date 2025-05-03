package auth_test

import (
	"errors"
	"fmt"
	"link-manager/internal/auth"
	"link-manager/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (mockRepo *MockUserRepository) Create(user *user.User) (*user.User, error) {
	return user, nil
}

func (mockRepo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, errors.New("record not found")
}

func TestRegisterSuccess(t *testing.T) {
	authService := auth.NewAuthService(&MockUserRepository{})
	user, err := authService.Register("name@domain.ru", "passw", "Test User")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user.Email, user.Name)
}
