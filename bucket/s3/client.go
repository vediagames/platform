package s3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/vediagames/vediagames.com/bucket/domain"
)

type client struct {
	client *s3manager.Uploader
	bucket string
}

type Config struct {
	Key      string
	Secret   string
	Region   string
	EndPoint string
	Bucket   string
}

func (c Config) Validate() error {
	if c.Key == "" {
		return fmt.Errorf("key is required")
	}

	if c.Secret == "" {
		return fmt.Errorf("secret is required")
	}

	if c.Region == "" {
		return fmt.Errorf("region is required")
	}

	if c.EndPoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	if c.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}

	return nil
}

func New(cfg Config) (domain.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.Key, cfg.Secret, ""),
		Endpoint:         aws.String(cfg.EndPoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         aws.LogLevel(aws.LogDebugWithHTTPBody),
	}

	s, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	c := s3manager.NewUploader(s)

	return &client{client: c, bucket: cfg.Bucket}, nil
}

func (s client) Upload(ctx context.Context, path string, reader io.Reader) error {
	contentType := "image/jpeg"

	if strings.Contains(path, ".svg") {
		contentType = "image/svg+xml"
	}

	_, err := s.client.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        reader,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})

	return fmt.Errorf("failed to upload: %w", err)
}
