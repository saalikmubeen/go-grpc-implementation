package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable
// and stored in the Config struct.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBURI                string        `mapstructure:"DB_URI"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	JWTSecret            string        `mapstructure:"JWT_SECRET"`
}

// LoadConfig reads configuration from configuration or environment files
// inside the given directory path if it exists or override their values
// with environment variables if they are provided. The values are
// then stored and loaded into the Config struct.
func LoadConfig(path string) (config Config, err error) {

	// Set the name of the config file to look for and
	// the path to look for the config file in:
	viper.AddConfigPath(path)  // The path to look for the config file in
	viper.SetConfigName("app") // The name of the config file without the extension
	viper.SetConfigType("env") // The extension of the config file: env, json, xml, yml, etc

	// To read values from environment variables:
	// Read values from environment variables and automatically override the
	// values read from the config file wih the values read from the
	// environment variables if they exist.
	viper.AutomaticEnv()

	// To read values from a config file:
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the values read into the Config struct:x
	err = viper.Unmarshal(&config)

	return config, err
}
