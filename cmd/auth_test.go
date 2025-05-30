package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"link-manager/internal/auth"
	"link-manager/internal/user"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func addUser(user *user.User) (*http.Response, error) {
	ts := httptest.NewServer(app(".env"))
	defer ts.Close()

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestRegisterSuccess(t *testing.T) {
	db := initDb()
	newUser := &user.User{
		Email:    "name@domian.ru",
		Password: "password",
		Name:     "Test User",
	}

	resp, err := addUser(newUser)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want status code: %d, got status code^ %d", http.StatusCreated, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	var regResp auth.RegisterResponse
	err = json.Unmarshal(data, &regResp)
	if err != nil {
		t.Fatal(err)
		return
	}

	var encryptedPass string

	db.Table("users").
		Select("password").
		Where("email = ? AND name = ?", newUser.Email, newUser.Name).
		Scan(&encryptedPass)

	fmt.Println(encryptedPass)
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "name@domain.ru",
		Password: "$2a$10$i5DPWWDE1CvbvTIw69leJeki34Sk8q4EDicdhaYpwFF8DYIeu0Ea.",
		Name:     "Test Name",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().Where("email = ?", "name@domain.ru").Delete(&user.User{})
}

func signIn(account *auth.LoginRequest) (*http.Response, error) {
	ts := httptest.NewServer(app(".env"))
	defer ts.Close()

	data, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestLoginSuccess(t *testing.T) {
	want := http.StatusOK

	account := &auth.LoginRequest{
		Email:    "name@domain.ru",
		Password: "pass",
	}

	db := initDb()
	initData(db)

	resp, err := signIn(account)
	if err != nil {
		t.Fatal(err)
	}

	got := resp.StatusCode

	if resp.StatusCode != want {
		t.Fatalf("want: %d, got: %d", want, got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var loginResp auth.LoginResponse
	json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		if resp.StatusCode != want {
			t.Fatalf("JWT-token empty")
		}
	}

	removeData(db)
}

func TestLoginWrongPass(t *testing.T) {
	want := http.StatusUnauthorized

	account := &auth.LoginRequest{
		Email:    "name@domain.ru",
		Password: "passa",
	}

	db := initDb()
	initData(db)

	resp, err := signIn(account)
	if err != nil {
		t.Fatal(err)
	}

	got := resp.StatusCode

	if resp.StatusCode != want {
		t.Fatalf("want: %d, got: %d", want, got)
	}

	removeData(db)
}

func TestLoginWrongUser(t *testing.T) {
	want := http.StatusUnauthorized

	account := &auth.LoginRequest{
		Email:    "nv.luntsev@yandex.ru",
		Password: "pass",
	}

	db := initDb()
	initData(db)

	resp, err := signIn(account)
	if err != nil {
		t.Fatal(err)
	}

	got := resp.StatusCode

	if resp.StatusCode != want {
		t.Fatalf("want: %d, got: %d", want, got)
	}

	removeData(db)
}
