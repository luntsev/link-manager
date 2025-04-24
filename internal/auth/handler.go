package auth

import (
	"fmt"
	"link-manager/configs"
	"link-manager/pkg/request"
	"link-manager/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		_, err := request.HandleBody[LoginRequest](&w, req)
		if err != nil {
			return
		}

		resp := LoginResponse{
			Token: "123",
		}

		response.Json(w, resp, 200)
		fmt.Println("Login")
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		body, err := request.HandleBody[RegisterRequest](&w, req)
		if err != nil {
			return
		}
		user, _ := handler.AuthService.Register(body.Email, body.Password, body.Name)

		response.Json(w, user, http.StatusOK) // 200
	}
}

func NewAuthHendler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}
