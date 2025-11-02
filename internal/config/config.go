package config

type AppConfig struct {
	App struct {
		Port      string `mapstructure:"port"`
		JwtSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"app"`
}
