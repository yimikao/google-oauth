package utils

import "github.com/spf13/viper"

type Config struct {
	GoogleOauthClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOauthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	CallbackURL             string `mapstructure:"CALLBACK_URL"`
}

func LoadConfig(path string) (*Config, error) {
	var (
		config Config
		err    error
	)
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return &config, err
	}

	err = viper.Unmarshal(&config)
	return &config, err
}
