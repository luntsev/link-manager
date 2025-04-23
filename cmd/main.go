package main

import (
	"fmt"
	"link-manager/configs"
	"link-manager/internal/auth"
	"link-manager/internal/link"
	"link-manager/internal/user"
	"link-manager/pkg/db"
	"link-manager/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	dataBase := db.NewDb(conf)
	router := http.NewServeMux()

	//Repository

	linkRepository := link.NewLinkRepository(dataBase)
	userRepository := user.NewUserRepository(dataBase)

	//Services
	authService := auth.NewAuthService(userRepository)

	//Hendlers
	auth.NewAuthHendler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHendler(router, link.LinkHandlerDeps{
		LinkRepo: linkRepository,
	})

	//Middlewares

	stackMw := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stackMw(router),
	}

	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()
}
