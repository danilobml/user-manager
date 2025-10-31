package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/danilobml/user-manager/internal/config"
	"github.com/danilobml/user-manager/internal/httpx"
	"github.com/danilobml/user-manager/internal/routes"
	"github.com/danilobml/user-manager/internal/user/handler"
	"github.com/danilobml/user-manager/internal/user/repository"
	"github.com/danilobml/user-manager/internal/user/service"
)

func main() {
	// config loading
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error reading config file, %s", err)
    }
	var config config.AppConfig
	err = viper.Unmarshal(&config)
    if err != nil {
        log.Fatalf("Unable to decode into struct, %v", err)
    }

	// Initializations
	userRepository := repository.NewUserRepositoryDdb()
	userService := service.NewUserserviceImpl(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := routes.NewRouter(userHandler)

	httpx.Serve(config.App.Port, &router)
}
