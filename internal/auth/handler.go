package auth

import (
	"link-manager/configs"
	"link-manager/pkg/jwt"
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

func NewAuthHendler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized) //401
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: body.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}

		resp := &LoginResponse{
			Token: token,
		}

		response.Json(w, resp, http.StatusOK) // 200
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := request.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: body.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}

		resp := &RegisterResponse{
			Token: token,
		}

		response.Json(w, resp, http.StatusOK) // 200
	}
}
