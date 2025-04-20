package main

import (
	"fmt"
	"link-manager/configs"
	"link-manager/internal/auth"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	auth.NewAuthHendler(router, auth.AuthHandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()
}
