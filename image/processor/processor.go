package processor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
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

func New(client http.Client, cfg config.Imagor) domain.Processor {
	return processor{
		client:    client,
		imagorURL: cfg.URL,
		secret:    cfg.Secret,
	}
}

func (p processor) Process(ctx context.Context, imageURL string) (io.Reader, error) {
	// TODO: Generate URL with thumbnail.go
	// TODO: Get format and resolution
	_ = p.generateUrl("", p.secret)
	res, err := http.Get("url")
	if err != nil {
		return nil, fmt.Errorf("failed to process the image: %w", err)
	}

	return res.Body, nil
}

func (p processor) generateUrl(path, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(path))
	s := base64.StdEncoding.EncodeToString(hash.Sum(nil))[:40]
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	return fmt.Sprintf("%s%s/%s", p.imagorURL, s, path)
}
