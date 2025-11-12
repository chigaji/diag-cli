package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) error {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("diag")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	//defaults
	viper.SetDefault("output", "table")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("process.default_top", 5)

	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	}

	return nil
}
