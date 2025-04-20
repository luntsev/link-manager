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
}

type AuthHandler struct {
	*configs.Config
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
		payload, err := request.HandleBody[RegisterRequest](&w, req)
		if err != nil {
			return
		}

		resp := *payload
		response.Json(w, resp, 200)
	}
}

func NewAuthHendler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}
