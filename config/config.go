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
	Auth struct {
		KratosURL string `mapstructure:"kratosURL"`
	} `mapstructure:"auth"`
	RedisAddress string       `mapstructure:"redisAddress"`
	Imagor       Imagor       `mapstructure:"imagor"`
	BunnyStorage BunnyStorage `mapstructure:"bunnyStorage"`
}

func (c Config) Validate() error {
	var err zeroerror.Error

	if c.Environment == "" {
		err.Add(fmt.Errorf("empty environment"))
	}

	if c.LogLevel == "" {
		err.Add(fmt.Errorf("empty log level"))
	}

	if c.Port == 0 {
		err.Add(fmt.Errorf("empty port"))
	}

	if c.PostgreSQL.ConnectionString == "" {
		err.Add(fmt.Errorf("empty postgresql connection string"))
	}

	if c.SendInBlue.Key == "" {
		err.Add(fmt.Errorf("empty sendinblue key"))
	}

	if c.Bucket.Key == "" {
		err.Add(fmt.Errorf("empty bucket key"))
	}

	if c.Bucket.Secret == "" {
		err.Add(fmt.Errorf("empty bucket secret"))
	}

	if c.Bucket.Region == "" {
		err.Add(fmt.Errorf("empty bucket region"))
	}

	if c.Bucket.EndPoint == "" {
		err.Add(fmt.Errorf("empty bucket endpoint"))
	}

	if c.Bucket.Bucket == "" {
		err.Add(fmt.Errorf("empty bucket name"))
	}

	if len(c.CORS.AllowedOrigins) == 0 {
		err.Add(fmt.Errorf("empty cors allowed origins"))
	}

	for _, origin := range c.CORS.AllowedOrigins {
		if origin == "" {
			err.Add(fmt.Errorf("empty cors allowed origin"))
		}
	}

	if c.PostgreSQL.Path.Migration == "" {
		err.Add(fmt.Errorf("empty postgresql path migration"))
	}

	if c.PostgreSQL.Path.Stub == "" {
		err.Add(fmt.Errorf("empty postgresql path stub"))
	}

	if c.RedisAddress == "" {
		err.Add(fmt.Errorf("empty redis address"))
	}

	if c.Auth.KratosURL == "" {
		err.Add(fmt.Errorf("empty kratos auth url"))
	}

	if c.BunnyStorage.URL == "" {
		err.Add(fmt.Errorf("empty bunny storage url key"))
	}

	if c.BunnyStorage.AccessKey == "" {
		err.Add(fmt.Errorf("empty bunny storage access key"))
	}

	if c.BunnyStorage.AccessKey == "" {
		err.Add(fmt.Errorf("empty bunny storage access key"))
	}

	if c.Imagor.URL == "" {
		err.Add(fmt.Errorf("empty imagor url key"))
	}

	if c.Imagor.Secret == "" {
		err.Add(fmt.Errorf("empty imagor secret key"))
	}

	if c.Imagor.S3CDNURL == "" {
		err.Add(fmt.Errorf("empty s3 cdn url key"))
	}

	if c.Imagor.BunnyCDNURL == "" {
		err.Add(fmt.Errorf("empty bunny cdn url key"))
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
