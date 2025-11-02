package main

import (
	"strings"

	"github.com/danilobml/user-manager/internal/config"
	"github.com/danilobml/user-manager/internal/httpx/middleware"

	//"github.com/danilobml/user-manager/internal/ses"
	// "github.com/danilobml/user-manager/internal/ddb"
	"github.com/danilobml/user-manager/internal/httpx"
	mail_service "github.com/danilobml/user-manager/internal/mailer/service"
	"github.com/danilobml/user-manager/internal/routes"
	user_handler "github.com/danilobml/user-manager/internal/user/handler"
	"github.com/danilobml/user-manager/internal/user/jwt"
	user_repository "github.com/danilobml/user-manager/internal/user/repository"
	user_service "github.com/danilobml/user-manager/internal/user/service"
)

func main() {
	config := config.LoadConfig()

	jwtManager := jwt.NewJwtManager([]byte(config.App.JwtSecret))
	// TODO: reactivate ddbRepo when infra implemented and deployed to AWS
	// ddbClient := ddb.InitDynamo()
	// userRepository := repository.NewUserRepositoryDdb(ddbClient)

	userRepository := user_repository.NewUserRepositoryInMemory()

	// TODO replace with SES mailer when infra implemented and deployed to AWS
	// sesMailClient := ses.Ses_Init()
	// mailService := mail_service.NewSesMailService(sesMailClient, config.Mail.FromEmail)

	// Local mailer
	mailService := mail_service.NewLocalMailService(mail_service.LocalMailConfig{
		FromEmail:     config.Mail.FromEmail,
		FromEmailPass: config.Mail.FromEmailPass,
		FromEmailSMTP: config.Mail.FromEmailSMTP,
		SMTPAddr:      config.Mail.SMTPAddr,
	})

	userService := user_service.NewUserserviceImpl(userRepository, jwtManager, mailService, config.App.BaseUrl)
	userHandler := user_handler.NewUserHandler(userService, strings.TrimSpace(config.App.ApiKey))

	authMiddleware := middleware.Authenticate(jwtManager)

	router := routes.NewRouter(userHandler, authMiddleware)

	httpx.Serve(config.App.Port, &router)
}
