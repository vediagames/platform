package config

import (
	"fmt"
	"log"

	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/vediagames/zeroerror"
)

type contextKey string

const (
	ContextKey contextKey = "config_context_key"
)

type Imagor struct {
	URL    string `mapstructure:"URL" validate:"required"`
	Secret string `mapstructure:"secret" validate:"required"`
}
type BunnyStorage struct {
	URL       string `mapstructure:"URL" validate:"required"`
	Zone      string `mapstructure:"zone" validate:"required"`
	AccessKey string `mapstructure:"accessKey" validate:"required"`
}

type Config struct {
	Environment string `mapstructure:"environment" validate:"required"`
	LogLevel    string `mapstructure:"logLevel" validate:"required"`
	Port        int    `mapstructure:"port" validate:"required"`
	PostgreSQL  struct {
		ConnectionString string `mapstructure:"connectionString" validate:"required"`
		Path             struct {
			Migration string `mapstructure:"migration" validate:"required"`
			Stub      string `mapstructure:"stub" validate:"required"`
		} `mapstructure:"path" validate:"required"`
	} `mapstructure:"postgresql" validate:"required"`
	SendInBlue struct {
		Key string `mapstructure:"key" validate:"required"`
	} `mapstructure:"sendinblue" validate:"required"`
	Bucket struct {
		Key      string `mapstructure:"key" validate:"required"`
		Secret   string `mapstructure:"secret" validate:"required"`
		Region   string `mapstructure:"region" validate:"required"`
		EndPoint string `mapstructure:"endpoint" validate:"required"`
		Bucket   string `mapstructure:"bucket" validate:"required"`
	} `mapstructure:"bucket" validate:"required"`
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowedOrigins" validate:"required"`
	} `mapstructure:"cors" validate:"required"`
	RedisAddress string       `mapstructure:"redisAddress" validate:"required"`
	Imagor       Imagor       `mapstructure:"imagor" validate:"required"`
	BunnyStorage BunnyStorage `mapstructure:"bunnyStorage" validate:"required"`
}

func (c Config) Validate() error {
	var err zeroerror.Error

	if c.Environment == "" {
		err.Add(fmt.Errorf("environment is not set"))
	}

	if c.LogLevel == "" {
		err.Add(fmt.Errorf("log level is not set"))
	}

	if c.Port == 0 {
		err.Add(fmt.Errorf("port is not set"))
	}

	if c.PostgreSQL.ConnectionString == "" {
		err.Add(fmt.Errorf("postgresql connection string is not set"))
	}

	if c.SendInBlue.Key == "" {
		err.Add(fmt.Errorf("sendinblue key is not set"))
	}

	if c.Bucket.Key == "" {
		err.Add(fmt.Errorf("bucket key is not set"))
	}

	if c.Bucket.Secret == "" {
		err.Add(fmt.Errorf("bucket secret is not set"))
	}

	if c.Bucket.Region == "" {
		err.Add(fmt.Errorf("bucket region is not set"))
	}

	if c.Bucket.EndPoint == "" {
		err.Add(fmt.Errorf("bucket endpoint is not set"))
	}

	if c.Bucket.Bucket == "" {
		err.Add(fmt.Errorf("bucket name is not set"))
	}

	if len(c.CORS.AllowedOrigins) == 0 {
		err.Add(fmt.Errorf("cors allowed origins is not set"))
	}

	for _, origin := range c.CORS.AllowedOrigins {
		if origin == "" {
			err.Add(fmt.Errorf("cors allowed origin is empty"))
		}
	}

	if c.PostgreSQL.Path.Migration == "" {
		err.Add(fmt.Errorf("postgresql path migration is not set"))
	}

	if c.PostgreSQL.Path.Stub == "" {
		err.Add(fmt.Errorf("postgresql path stub is not set"))
	}

	if c.RedisAddress == "" {
		err.Add(fmt.Errorf("redis address is not set"))
	}

	// TODO: add validation
	return err.Err()
}

func New(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var validate = validator.New()
	if err := validate.Struct(&config); err != nil {
		log.Fatalf("Missing required config variables %v\n", err)
	}

	return config, nil
}
