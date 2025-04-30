package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTData struct {
	Email string
}

type JWT struct {
	Secret string
}

func (j *JWT) Create(data JWTData) (string, error) {
	jToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})

	token, err := jToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"]
	return t.Valid, &JWTData{
		Email: email.(string),
	}
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}
