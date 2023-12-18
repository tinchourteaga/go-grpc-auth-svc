package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Load() error {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msg("error reading config file")
		return err
	}

	return nil
}
