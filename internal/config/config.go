package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type FCMConfig struct {
	CredentialsFile string   `mapstructure:"credentials_file"`
	Scopes          []string `mapstructure:"scopes"`
	EndpointURL     string   `mapstructure:"endpoint_url"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	FCM    FCMConfig    `mapstructure:"fcm"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
