package jwt_test

import (
	"link-manager/pkg/jwt"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	data := jwt.JWTData{
		Email: "name@domain.ru",
	}

	want := data.Email

	testJWT := jwt.NewJWT("4HAgYbhrfM6k33p7zQPOXknFoczQjqK3SMxHIgGm4DR5g8J9YzwzcsFWlniSpb")

	jwtToken, err := testJWT.Create(data)
	if err != nil {
		t.Fatal(err)
	}

	isValid, resData := testJWT.Parse(jwtToken)
	if !isValid {
		t.Fatal("jwt-token is not walid")
	}

	got := resData.Email

	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
