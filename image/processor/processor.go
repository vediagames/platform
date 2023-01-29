package processor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/vediagames/vediagames.com/config"
	"github.com/vediagames/vediagames.com/image/domain"
)

type processor struct {
	client    http.Client
	imagorURL string
	secret    string
}

type Config struct {
	HTTPClient http.Client
	ImagorCfg  config.Imagor
}

func New(cfg Config) domain.Processor {

	return processor{
		client:    cfg.HTTPClient,
		imagorURL: cfg.ImagorCfg.URL,
		secret:    cfg.ImagorCfg.Secret,
	}
}

func (p processor) Process(ctx context.Context, req domain.ProcessQuery) (domain.ProcessResult, error) {

	reqImageURL := fmt.Sprintf("fit-in/%dx%d/filters:format(%s)/%s", req.Thumbnail.Width, req.Thumbnail.Height, req.Thumbnail.Format, req.ImageURL)
	encryptedURL := p.generateUrl(reqImageURL, p.secret)

	res, err := p.client.Get(encryptedURL)
	if err != nil {
		return domain.ProcessResult{}, fmt.Errorf("failed to process the image: %w", err)
	}

	return domain.ProcessResult{
		File: res.Body,
	}, nil
}

func (p processor) generateUrl(path, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(path))

	s := base64.StdEncoding.EncodeToString(hash.Sum(nil))[:40]

	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)

	return fmt.Sprintf("%s%s/%s", p.imagorURL, s, path)
}
