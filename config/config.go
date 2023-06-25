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
	S3 struct {
		Key      string `mapstructure:"key"`
		Secret   string `mapstructure:"secret"`
		Region   string `mapstructure:"region"`
		Endpoint string `mapstructure:"endpoint"`
		Bucket   string `mapstructure:"bucket"`
	}
	QuotesCSV string `mapstructure:"quotesCSV"`
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Environment == "", fmt.Errorf("environment is not set"))
	err.AddIf(c.LogLevel == "", fmt.Errorf("logLevel is not set"))
	err.AddIf(c.Port == 0, fmt.Errorf("port is not set"))
	err.AddIf(c.PostgreSQL.ConnectionString == "", fmt.Errorf("postgresql.connectionString is not set"))
	err.AddIf(c.SendInBlue.Key == "", fmt.Errorf("sendinblue.key is not set"))
	err.AddIf(len(c.CORS.AllowedOrigins) == 0, fmt.Errorf("cors.allowedOrigins is not set"))
	err.AddIf(c.PostgreSQL.Path.Migration == "", fmt.Errorf("postgresql.path.migration is not set"))
	err.AddIf(c.PostgreSQL.Path.Stub == "", fmt.Errorf("postgresql.path.stub is not set"))
	err.AddIf(c.RedisAddress == "", fmt.Errorf("redisAddress is not set"))
	err.AddIf(c.Auth.KratosURL == "", fmt.Errorf("auth.kratusURL is not set"))
	err.AddIf(c.BigQuery.ProjectID == "", fmt.Errorf("bigquery.projectID is not set"))
	err.AddIf(c.BigQuery.CredentialsPath == "", fmt.Errorf("bigquery.credentialsPath is not set"))
	err.AddIf(c.Imagor.URL == "", fmt.Errorf("imagor.URL is not set"))
	err.AddIf(c.Imagor.Secret == "", fmt.Errorf("imagor.secret is not set"))
	err.AddIf(c.S3.Key == "", fmt.Errorf("s3.key is not set"))
	err.AddIf(c.S3.Secret == "", fmt.Errorf("s3.secret is not set"))
	err.AddIf(c.S3.Endpoint == "", fmt.Errorf("s3.endpoint is not set"))
	err.AddIf(c.S3.Bucket == "", fmt.Errorf("s3.bucket is not set"))
	err.AddIf(c.QuotesCSV == "", fmt.Errorf("quotesCSV is not set"))

	for _, origin := range c.CORS.AllowedOrigins {
		err.AddIf(origin == "", fmt.Errorf("cors.allowedOrigins includes empty origin"))
	}

	return err.Err()
}

func New(path string) Config {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return config
}
