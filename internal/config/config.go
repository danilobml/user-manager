package config

type AppConfig struct {
    App struct {
        Port string    `mapstructure:"port"`
    } `mapstructure:"app"`
}
