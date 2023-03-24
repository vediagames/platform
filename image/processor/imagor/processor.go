package imagor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/vediagames/zeroerror"

	bucketdomain "github.com/vediagames/platform/bucket/domain"
	"github.com/vediagames/platform/image/domain"
)

type Config struct {
	URL          string
	Secret       string
	Client       *http.Client
	BucketClient bucketdomain.Client
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.Client == nil, fmt.Errorf("empty client"))
	err.AddIf(c.BucketClient == nil, fmt.Errorf("empty bucket client"))
	err.AddIf(c.URL == "", fmt.Errorf("empty URL"))
	err.AddIf(c.Secret == "", fmt.Errorf("empty secret"))

	return err.Err()
}

func New(c Config) domain.Processor {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return &processor{
		url:          c.URL,
		secret:       c.Secret,
		client:       c.Client,
		bucketClient: c.BucketClient,
	}
}

type processor struct {
	url          string
	secret       string
	client       *http.Client
	bucketClient bucketdomain.Client
}

func (p processor) Process(ctx context.Context, req domain.ProcessRequest) (domain.ProcessResponse, error) {
	reqImageURL := fmt.Sprintf("stretch/%dx%d/filters:format(%s)/%s",
		req.Image.Width,
		req.Image.Height,
		req.Image.Format,
		req.OriginalImageURL,
	)

	hash := hmac.New(sha256.New, []byte(p.secret))
	hash.Write([]byte(reqImageURL))
	s := base64.StdEncoding.EncodeToString(hash.Sum(nil))[:40]
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)

	encryptedURL := fmt.Sprintf("%s/%s/%s", p.url, s, reqImageURL)

	httpRes, err := p.client.Get(encryptedURL)
	if err != nil {
		return domain.ProcessResponse{}, fmt.Errorf("failed to get: %w", err)
	}

	if err := p.bucketClient.Upload(ctx, req.Path, httpRes.Body); err != nil {
		return domain.ProcessResponse{}, fmt.Errorf("failed to upload: %w", err)
	}

	httpRes.Body.Close()

	return domain.ProcessResponse{
		Path: req.Path,
	}, nil
}
