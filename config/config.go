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

type Imagor struct {
	URL         string `mapstructure:"URL"`
	Secret      string `mapstructure:"secret"`
	S3CDNURL    string `mapstructure:"s3CDNURL"`
	BunnyCDNURL string `mapstructure:"bunnyCDNURL"`
}
type BunnyStorage struct {
	URL       string `mapstructure:"URL"`
	Zone      string `mapstructure:"zone"`
	AccessKey string `mapstructure:"accessKey"`
}

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
		Bucket   string `mapstructure:"bucket"`
	} `mapstructure:"bucket"`
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowedOrigins"`
	} `mapstructure:"cors"`
	RedisAddress string       `mapstructure:"redisAddress"`
	Imagor       Imagor       `mapstructure:"imagor"`
	BunnyStorage BunnyStorage `mapstructure:"bunnyStorage"`
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

	if c.BunnyStorage.URL == "" {
		err.Add(fmt.Errorf("bunny storage url key is not set"))
	}

	if c.BunnyStorage.AccessKey == "" {
		err.Add(fmt.Errorf("bunny storage access key is not set"))
	}

	if c.BunnyStorage.AccessKey == "" {
		err.Add(fmt.Errorf("bunny storage access key is not set"))
	}

	if c.Imagor.URL == "" {
		err.Add(fmt.Errorf("imagor url key is not set"))
	}

	if c.Imagor.Secret == "" {
		err.Add(fmt.Errorf("imagor secret key is not set"))
	}

	if c.Imagor.S3CDNURL == "" {
		err.Add(fmt.Errorf("s3 cdn url key is not set"))
	}

	if c.Imagor.BunnyCDNURL == "" {
		err.Add(fmt.Errorf("bunny cdn url key is not set"))
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
