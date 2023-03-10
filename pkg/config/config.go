// Package config - Configuration
package config

import (
	"os"
	"strings"

	"kore-on/pkg/logger"

	"github.com/spf13/viper"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// Load - Load the configuration from file
func Load() error {
	var err error

	// dir := strings.Split(workingdir, "/")
	path := "./conf"
	// Search config files in config directory with name "config.yaml" (without extension).
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// yaml의 속성 ex) database.database_name 등을 update 위한 환경변수 설정시 . 대신 '_' 사용할 수 있게 한다.
	// ex)  database.database_name 속성은  DATABASE_DATABASE_NAME 환경변수 설정하면 값이 override 됨
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logger.Errorf("Could not load configuration file: %s", err.Error())
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
		}
	}

	logger.Info("Using config file:", viper.ConfigFileUsed())

	return err
}
