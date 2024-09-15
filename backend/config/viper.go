package config

import (
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	AUTH0_DOMAIN   string `mapstructure:"AUTH0_DOMAIN"`
	AUTH0_AUDIENCE string `mapstructure:"AUTH0_AUDIENCE"`
	PORT           string `mapstructure:"PORT"`
	DATABASE_URL   string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (config EnvVars, err error) {
	viper.AutomaticEnv()

	if _, err := os.Stat(".env.local"); err == nil {
		viper.AddConfigPath(".")
		viper.SetConfigName(".env.local")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			return config, err
		}
	}

	err = viper.Unmarshal(&config)
	return
}
