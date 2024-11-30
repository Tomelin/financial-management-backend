package main

import (
	"log"

	"github.com/synera-br/financial-management/src/backend/configs"
	"github.com/synera-br/financial-management/src/backend/internal/core/repository"
	"github.com/synera-br/financial-management/src/backend/internal/core/service"
	"github.com/synera-br/financial-management/src/backend/internal/infra/handler/web"
	"github.com/synera-br/financial-management/src/backend/pkg/authProvider"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	http_server "github.com/synera-br/financial-management/src/backend/pkg/http_server/server"
)

func main() {

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fbDB, err := db.NewFirebaseDatabaseConnection(cfg.Fields["firebase"], "firebase")

	authProvider, err := authProvider.NewAuthProvider(cfg.Fields["auth"])
	if err != nil {
		log.Fatalln(err)
	}

	store, err := authProvider.Store()
	if err != nil {
		log.Fatalln(err)
	}

	rest, err := http_server.NewRestApi(cfg.Fields["webserver"], store)
	if err != nil {
		log.Fatalln(err)
	}

	// USER
	userRepo, err := repository.NewUserRepository(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	userSvc, err := service.NewUserService(userRepo)
	if err != nil {
		log.Fatalln(err)
	}

	// CATEGORY
	repoCategory, err := repository.NewCategoryRepository(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	svcCategory, err := service.NewCategorySvcsitory(repoCategory)
	if err != nil {
		log.Fatalln(err)
	}

	// WEbServer
	// validate := web.NewAuthHandlerHttp(authProvider, userSvc, rest.RouterGroup)
	web.NewUserHandlerHttp(&userSvc, rest.RouterGroup)
	web.NewCategoryHandlerHttp(svcCategory, rest.RouterGroup)

	rest.Run(rest.Route.Handler())
}
