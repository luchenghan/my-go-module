package config

import (
	"fmt"
	"mymodule/gin/pkg/logger"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DataSourceName  string        `mapstructure:"dataSourceName"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetㄊime"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expiry int    `mapstructure:"expiry"`
}

// Global config variable
var GlobalConfig *Config

// LoadConfig initializes and reads the configuration
func LoadConfig() error {
	// Set config file details
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // config file type
	viper.AddConfigPath("./config") // path to look for the config file

	// Set default values (optional)
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("jwt.expiry", 3600) // default expiry time in 1 hour

	// Allow environment variables to override config
	viper.AutomaticEnv()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal into struct
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(GlobalConfig); err != nil {
			fmt.Printf("Error reloading config: %v\n", err)
		}

		logger.Initialize(GlobalConfig.Logger.Encoding, GlobalConfig.Logger.Level)
	})

	return nil
}
