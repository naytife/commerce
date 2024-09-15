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
	viper.BindEnv("AUTH0_DOMAIN")   // Binds $AUTH0_DOMAIN to viper.Get("AUTH0_DOMAIN")
	viper.BindEnv("AUTH0_AUDIENCE") // Binds $AUTH0_AUDIENCE to viper.Get("AUTH0_AUDIENCE")
	viper.BindEnv("PORT")           // Binds $PORT to viper.Get("PORT")
	viper.BindEnv("DATABASE_URL")

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
