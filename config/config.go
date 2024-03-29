package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config wraps the config file as a struct
type Config struct {
	Version     string
	LogLevel    log.Level
	JwtSecret   string
	AdminSecret string
	API         apiConfig
	Database    databaseConfig
}

type apiConfig struct {
	Port      int
	SecretKey string
}

type databaseConfig struct {
	Host     string
	Port     int
	Database string
	User     string
}

// init serializes YAML into a Config struct
func (cfg *Config) init() {
	cfg.Version = viper.GetString("version")
	cfg.setLogLevel(viper.GetString("loglevel"))
	cfg.JwtSecret = viper.GetString("jwt_secret")
	cfg.AdminSecret = viper.GetString("admin_secret")
	cfg.API.Port = viper.GetInt("api.port")
	cfg.API.SecretKey = viper.GetString("api.secret_key")
	cfg.Database.Host = viper.GetString("db.host")
	cfg.Database.Port = viper.GetInt("db.port")
	cfg.Database.Database = viper.GetString("db.database")
	cfg.Database.User = viper.GetString("db.user")
}

// GetConfig loads config data into a Config struct
func GetConfig() *Config {
	cfg := new(Config)
	cfg.init()

	return cfg
}

// InitConfig sets up the config file
func InitConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	cfg := GetConfig()
	log.SetLevel(cfg.LogLevel)

	return cfg, nil
}

func (cfg *Config) setLogLevel(loglevel string) {
	switch loglevel {
	case "info":
		cfg.LogLevel = log.InfoLevel
	case "warn":
		cfg.LogLevel = log.WarnLevel
	case "fatal":
		cfg.LogLevel = log.FatalLevel
	}
}
