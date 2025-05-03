package auth_test

import (
	"link-manager/internal/user"
	"link-manager/pkg/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginSuccess(t *testing.T) {
	dataBase, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed init mock DB")
		return
	}

	mockDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dataBase,
	}))
	if err != nil {
		t.Fatalf("Failed init GORM")
		return
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: mockDb,
	})

}
