package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port     string
	Env      string
	LogLevel string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int
}func LoadConfig() Config {
	viper.Setconfigfile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Waring: Failed to read the config file: $s\n", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port:	 viper.GetString("SERVER_PORT"),
			Env:    viper.GetString("SERVER_ENV"),
			LogLevel: viper.GetString("SERVER_LOG_LEVEL"),
		},
		Database: DatabaseConfig{
			Host: 		  viper.GetString("DB_HOST"),
			Port: 		  viper.GetString("DB_PORT"),
			User: 		  viper.GetString("DB_USER"),
			Password: 	  viper.GetString("DB_PASSWORD"),
			Name: 		  viper.GetString("DB_NAME"),
			Charset: 	  viper.GetString("DB_CHARSET"),
			MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			SecretKey: viper.GetString("JWT_SECRET"),
			ExpiresIn: viper.GetInt("JWT_EXPIRES_IN"),
		},
	}

	//set default value
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Env == "" {
		config.Server.Env = "development"
	}
	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "debug"
	}
	if config.Database.Charset == "" {
		config.Database.Charset = "utf8mb4"
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = 10 * time.Hour
	}
	if config.JWT.ExpiresIn == 0 {
		config.JWT.ExpiresIn = "24h"
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset)
}