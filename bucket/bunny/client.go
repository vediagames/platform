package bunny

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vediagames/zeroerror"

	"github.com/vediagames/platform/bucket/domain"
)

type client struct {
	url       string
	accessKey string
	zone      string
	client    *http.Client
}

type Config struct {
	URL       string
	AccessKey string
	Zone      string
	Client    *http.Client
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.URL == "", fmt.Errorf("empty URL"))
	err.AddIf(c.AccessKey == "", fmt.Errorf("empty access key"))
	err.AddIf(c.Zone == "", fmt.Errorf("empty zone"))
	err.AddIf(c.Client == nil, fmt.Errorf("empty client"))

	return err.Err()
}

func New(c Config) domain.Client {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &client{
		url:       c.URL,
		accessKey: c.AccessKey,
		zone:      c.Zone,
		client:    c.Client,
	}
}

func (s client) Upload(ctx context.Context, path string, reader io.Reader) error {
	url := fmt.Sprintf("%s/%s/%s", s.url, s.zone, path)

	req, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("AccessKey", s.accessKey)

	switch {
	case strings.Contains(path, ".webp"):
		req.Header.Add("content-type", "image/webp")
	case strings.Contains(path, ".png"):
		req.Header.Add("content-type", "image/png")
	case strings.Contains(path, ".jpg"):
		req.Header.Add("content-type", "image/jpeg")
	default:
		return fmt.Errorf("failed to find appropriate content type")
	}

	httpRes, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do: %w", err)
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode > 299 {
		bytes, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return fmt.Errorf("failed to read all: %w", err)
		}

		return fmt.Errorf("failed: %s", string(bytes))
	}

	return nil
}
