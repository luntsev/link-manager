package auth_test

import (
	"bytes"
	"encoding/json"
	"link-manager/configs"
	"link-manager/internal/auth"
	"link-manager/internal/user"
	"link-manager/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	dataBase, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	mockDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dataBase,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: mockDb,
	})

	handler := auth.AuthHandler{
		AuthService: auth.NewAuthService(userRepo),
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
	}
	return &handler, mock, nil
}

func TestLoginRegisterSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}

	rows := mock.NewRows([]string{"created_at", "updated_at", "deleted_at", "email", "password", "Name"})
	//AddRow("name@domain.ru", "$2a$10$i5DPWWDE1CvbvTIw69leJeki34Sk8q4EDicdhaYpwFF8DYIeu0Ea.", "Test Name")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, err := json.Marshal(&auth.RegisterRequest{
		Name: "Test User",
		LoginRequest: auth.LoginRequest{
			Email:    "name@domain.ru",
			Password: "pass",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.UserRegister()(w, r)
	gotStatusCode := w.Result().StatusCode
	wantStatusCode := http.StatusCreated
	if gotStatusCode != wantStatusCode {
		t.Fatalf("want: %d, got: %d", wantStatusCode, gotStatusCode)
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}

	rows := mock.NewRows([]string{"email", "password"}).AddRow("name@domain.ru", "$2a$10$i5DPWWDE1CvbvTIw69leJeki34Sk8q4EDicdhaYpwFF8DYIeu0Ea.")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	data, err := json.Marshal(&auth.LoginRequest{
		Email:    "name@domain.ru",
		Password: "pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.UserLogin()(w, r)
	gotStatusCode := w.Result().StatusCode
	wantStatusCode := http.StatusOK
	if gotStatusCode != wantStatusCode {
		t.Fatalf("want: %d, got: %d", wantStatusCode, gotStatusCode)
	}
}
