package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"strings"

	"github.com/danilobml/user-manager/internal/config"
	"github.com/danilobml/user-manager/internal/ddb"
	"github.com/danilobml/user-manager/internal/httpx/middleware"
	mail_service "github.com/danilobml/user-manager/internal/mailer/service"
	"github.com/danilobml/user-manager/internal/routes"
	"github.com/danilobml/user-manager/internal/ses"
	user_handler "github.com/danilobml/user-manager/internal/user/handler"
	"github.com/danilobml/user-manager/internal/user/jwt"
	user_repository "github.com/danilobml/user-manager/internal/user/repository"
	user_service "github.com/danilobml/user-manager/internal/user/service"
)

func buildHandler() *httpadapter.HandlerAdapter {
	cfg := config.LoadConfig()

	jwtManager := jwt.NewJwtManager([]byte(cfg.App.JwtSecret))
	ddbClient := ddb.InitDynamo()
	userRepository := user_repository.NewUserRepositoryDdb(ddbClient)

	sesMailClient := ses.Ses_Init()
	mailService := mail_service.NewSesMailService(sesMailClient, cfg.Mail.FromEmail)

	userService := user_service.NewUserserviceImpl(userRepository, jwtManager, mailService, cfg.App.BaseUrl)
	userHandler := user_handler.NewUserHandler(userService, strings.TrimSpace(cfg.App.ApiKey))

	authMiddleware := middleware.Authenticate(jwtManager)
	router := routes.NewRouter(userHandler, authMiddleware)

	return httpadapter.New(router)
}

func main() {
	adapter := buildHandler()

	lambda.Start(adapter.ProxyWithContext)
}
