package auth

type LoginRequest struct {
	Login    string `json:"login" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name string `json:"name" validate:"required"`
	LoginRequest
}
