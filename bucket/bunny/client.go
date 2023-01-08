package bunny

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/vediagames/vediagames.com/bucket/domain"
	"github.com/vediagames/vediagames.com/config"
)

type client struct {
	client    http.Client
	url       string
	accessKey string
	zone      string
}

type Config struct {
	HTTPClient http.Client
	BunnyCfg   config.BunnyStorage
}

func (c Config) Validate() error {
	if c.BunnyCfg.URL == "" {
		return fmt.Errorf("url is required")
	}

	if c.BunnyCfg.AccessKey == "" {
		return fmt.Errorf("accesskey is required")
	}

	return nil
}

func New(cfg Config) (domain.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &client{
		client:    cfg.HTTPClient,
		accessKey: cfg.BunnyCfg.AccessKey,
		url:       cfg.BunnyCfg.URL,
	}, nil
}

func (s client) Upload(ctx context.Context, path string, reader io.Reader) error {
	// baseCDNURL/{storageZoneName}/{path}/{fileName}
	url := fmt.Sprintf("/%s/%s/%s", s.url, s.zone, path)
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return fmt.Errorf("invalid http request to Bunny, %w", err)
	}
	req.Header.Add("AccessKey", s.accessKey)
	// TODO: generate content-type according to requested format
	req.Header.Add("content-type", "application/octet-stream")

	_, err = s.client.Do(req)
	return err
}
