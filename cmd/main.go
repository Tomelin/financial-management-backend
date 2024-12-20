package main

import (
	"context"
	"log"

	"github.com/synera-br/financial-management/src/backend/configs"
	"github.com/synera-br/financial-management/src/backend/internal/core/repository"
	"github.com/synera-br/financial-management/src/backend/internal/core/service"
	"github.com/synera-br/financial-management/src/backend/internal/infra/handler/web"
	"github.com/synera-br/financial-management/src/backend/pkg/authProvider"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	http_server "github.com/synera-br/financial-management/src/backend/pkg/http_server/server"
	"github.com/synera-br/financial-management/src/backend/pkg/logger"
	"github.com/synera-br/financial-management/src/backend/pkg/observability"
)

func main() {

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	customLogger := logger.NewLoggerConfig(cfg.Fields["logger"])

	tracer, cleanup, err := observability.InicializeTracer(cfg.Fields["otel"])
	if err != nil {
		log.Fatalln(err)

	}
	defer cleanup()

	ctx, span := tracer.Trace.Start(context.Background(), "main")
	defer span.End()

	fbDB, err := db.NewFirebaseDatabaseConnection(ctx, cfg.Fields["firebase"], "firebase")
	if err != nil {
		log.Fatalln(err)
	}

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

	// Plan
	repoPlan, err := repository.NewPlanRepository(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	svcPlan, err := service.NewPlanService(repoPlan)
	if err != nil {
		log.Fatalln(err)
	}

	// TENANT
	repoTenant, err := repository.NewTenantRepository(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	svcTenant, err := service.NewTenantService(repoTenant, svcPlan)
	if err != nil {
		log.Fatalln(err)
	}

	// USER
	userRepo, err := repository.NewUserRepository(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	userSvc, err := service.NewUserService(userRepo, svcTenant, tracer)
	if err != nil {
		log.Fatalln(err)
	}

	// WALLET
	repoWallet, err := repository.NewWalletRepo(fbDB)
	if err != nil {
		log.Fatalln(err)
	}

	svcWallet, mErr := service.NewWalletSvc(repoWallet, svcTenant, userSvc)
	if err != nil {
		log.Fatalln(mErr)
	}

	// CATEGORY
	repoCategory, mErr := repository.NewTransactionCategoryRepo(tracer, fbDB)
	if mErr != nil {
		log.Fatalln(err)
	}

	svcCategory, mErr := service.NewTransactionCategorySvc(tracer, &repoCategory, &svcTenant, &svcWallet, &userSvc)
	if mErr != nil {
		log.Fatalln(err)
	}

	// Authorization
	repoAuth, err := repository.NewAuthorizationRepo(fbDB, customLogger)
	if err != nil {
		log.Fatalln(err)
	}

	svcAuth, err := service.NewAuthorizationSvc(repoAuth, customLogger)
	if err != nil {
		log.Fatalln(err)
	}

	// WEbServer
	web.NewAuthenticationHandlerHttp(authProvider, customLogger, svcAuth, userSvc, rest.RouterGroup)
	web.NewUserHandlerHttp(&userSvc, tracer, rest.RouterGroup)
	// web.NewCategoryHandlerHttp(&svcCategory, rest.RouterGroup)
	web.NewTenantHandlerHttp(&svcTenant, rest.RouterGroup)
	web.NewPlanHandlerHttp(&svcPlan, rest.RouterGroup)
	web.NewWalletHandlerHttp(&svcWallet, &userSvc, rest.RouterGroup)
	web.NewTransactionCategoryHandlerHttp(tracer, &svcCategory, &svcWallet, rest.RouterGroup, rest.MiddlewareHeader)
	rest.Run(rest.Route.Handler())
}

func second(l logger.Logger) {

	l.Error(&logger.Message{Body: "second", Code: logger.ResponseCodeAccepted})
}
