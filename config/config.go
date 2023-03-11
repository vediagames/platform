package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vediagames/zeroerror"
)

type contextKey string

const (
	ContextKey contextKey = "config_context_key"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"logLevel"`
	Port        int    `mapstructure:"port"`
	PostgreSQL  struct {
		ConnectionString string `mapstructure:"connectionString"`
		Path             struct {
			Migration string `mapstructure:"migration"`
			Stub      string `mapstructure:"stub"`
		} `mapstructure:"path"`
	} `mapstructure:"postgresql"`
	SendInBlue struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"sendinblue"`
	Bucket struct {
		Key      string `mapstructure:"key"`
		Secret   string `mapstructure:"secret"`
		Region   string `mapstructure:"region"`
		EndPoint string `mapstructure:"endpoint"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"bucket"`
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowedOrigins"`
	} `mapstructure:"cors"`
	Auth struct {
		KratosURL string `mapstructure:"kratosURL"`
	} `mapstructure:"auth"`
	RedisAddress string `mapstructure:"redisAddress"`
	BigQuery     struct {
		ProjectID       string `mapstructure:"projectID"`
		CredentialsPath string `mapstructure:"credentialsPath"`
	} `mapstructure:"bigQuery"`
	Imagor struct {
		URL    string `mapstructure:"URL"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"imagor"`
	BunnyStorage struct {
		URL       string `mapstructure:"URL"`
		Zone      string `mapstructure:"zone"`
		AccessKey string `mapstructure:"accessKey"`
	} `mapstructure:"bunnyStorage"`
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Environment == "", fmt.Errorf("environment is not set"))
	err.AddIf(c.LogLevel == "", fmt.Errorf("logLevel is not set"))
	err.AddIf(c.Port == 0, fmt.Errorf("port is not set"))
	err.AddIf(c.PostgreSQL.ConnectionString == "", fmt.Errorf("postgresql.connectionString is not set"))
	err.AddIf(c.SendInBlue.Key == "", fmt.Errorf("sendinblue.key is not set"))
	err.AddIf(c.Bucket.Key == "", fmt.Errorf("bucket.key is not set"))
	err.AddIf(c.Bucket.Secret == "", fmt.Errorf("bucket.secret is not set"))
	err.AddIf(c.Bucket.Region == "", fmt.Errorf("bucket.region is not set"))
	err.AddIf(c.Bucket.EndPoint == "", fmt.Errorf("bucket.endpoint is not set"))
	err.AddIf(c.Bucket.Name == "", fmt.Errorf("bucket.name is not set"))
	err.AddIf(len(c.CORS.AllowedOrigins) == 0, fmt.Errorf("cors.allowedOrigins is not set"))
	err.AddIf(c.PostgreSQL.Path.Migration == "", fmt.Errorf("postgresql.path.migration is not set"))
	err.AddIf(c.PostgreSQL.Path.Stub == "", fmt.Errorf("postgresql.path.stub is not set"))
	err.AddIf(c.RedisAddress == "", fmt.Errorf("redisAddress is not set"))
	err.AddIf(c.Auth.KratosURL == "", fmt.Errorf("auth.kratusURL is not set"))
	err.AddIf(c.BigQuery.ProjectID == "", fmt.Errorf("bigquery.projectID is not set"))
	err.AddIf(c.BigQuery.CredentialsPath == "", fmt.Errorf("bigquery.credentialsPath is not set"))
	err.AddIf(c.BunnyStorage.URL == "", fmt.Errorf("bunnyStorage url key is not set"))
	err.AddIf(c.BunnyStorage.AccessKey == "", fmt.Errorf("bunnyStorage.accessKey is not set"))
	err.AddIf(c.Imagor.URL == "", fmt.Errorf("imagor.URL is not set"))
	err.AddIf(c.Imagor.Secret == "", fmt.Errorf("imagor.secret is not set"))

	for _, origin := range c.CORS.AllowedOrigins {
		err.AddIf(origin == "", fmt.Errorf("cors.allowedOrigins includes empty origin"))
	}

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

	return config, nil
}
