package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("warning: could not read config file: %v (env may still provide values)", err)
	}
	viper.AutomaticEnv()
	_ = viper.BindEnv("app.base_url", "BASE_URL")
	_ = viper.BindEnv("app.api_key", "API_KEY")
	_ = viper.BindEnv("mail.from_email", "FROM_EMAIL")
	_ = viper.BindEnv("mail.from_email_password", "FROM_EMAIL_PASSWORD")
	_ = viper.BindEnv("mail.from_email_smtp", "FROM_EMAIL_SMTP")
	_ = viper.BindEnv("mail.smtp_addr", "SMTP_ADDR")

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return config
}
