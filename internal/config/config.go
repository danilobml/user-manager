package config

type AppConfig struct {
	App struct {
		Port      string `mapstructure:"port"`
		JwtSecret string `mapstructure:"jwt_secret"`
		BaseUrl   string `mapstructure:"base_url"`
	} `mapstructure:"app"`

	Mail struct {
		FromEmail        string `mapstructure:"from_email"`
		FromEmailPass    string `mapstructure:"from_email_password"`
		FromEmailSMTP    string `mapstructure:"from_email_smtp"`
		SMTPAddr         string `mapstructure:"smtp_addr"`
	} `mapstructure:"mail"`
}
