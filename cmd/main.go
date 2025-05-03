package main

import (
	"fmt"
	"link-manager/configs"
	"link-manager/internal/auth"
	"link-manager/internal/link"
	"link-manager/internal/stat"
	"link-manager/internal/user"
	"link-manager/pkg/db"
	"link-manager/pkg/event"
	"link-manager/pkg/middleware"
	"net/http"
)

func app(envFile string) http.Handler {
	conf := configs.LoadConfig(envFile)
	dataBase := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//Repository

	linkRepository := link.NewLinkRepository(dataBase)
	userRepository := user.NewUserRepository(dataBase)
	statRepository := stat.NewStatRepository(dataBase)

	//Services
	authService := auth.NewAuthService(userRepository)

	statService := stat.NewStatService(stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	//Hendlers
	auth.NewAuthHendler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHendler(router, link.LinkHandlerDeps{
		LinkRepo: linkRepository,
		EventBus: eventBus,
		Config:   conf,
	})

	stat.NewStatHendler(router, stat.StatHandlerDeps{
		Config:         conf,
		StatRepository: statRepository,
	})

	go statService.AddClick()

	//Middlewares
	stackMw := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stackMw(router)
}

func main() {

	stack := app("../.env")
	server := http.Server{
		Addr:    ":8081",
		Handler: stack,
	}

	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()
}
