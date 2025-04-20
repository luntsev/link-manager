package token

import "math/rand"

func GenToken(length int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUXYVWZabcdefghigklmnopqrstuxyvwz1234567890")
	token := make([]rune, length)

	for i := range token {
		token[i] = letterRunes[rand.Intn(len(letterRunes)-1)]
	}
	return string(token)
}
