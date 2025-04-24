package jwt

import "github.com/golang-jwt/jwt/v5"

type JWT struct {
	Secret string
}

func (j *JWT) Create(email string) (string, error) {
	jToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	token, err := jToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}
