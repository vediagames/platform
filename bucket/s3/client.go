package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vediagames/platform/bucket/domain"
)

type client struct {
	client *s3.Client
	bucket string
}

type Config struct {
	Key      string
	Secret   string
	Region   string
	Endpoint string
	Bucket   string
}

func (c Config) Validate() error {

	if c.Key == "" {
		return fmt.Errorf("key is required")
	}

	if c.Secret == "" {
		return fmt.Errorf("secret is required")
	}

	if c.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	if c.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}

	return nil
}

func New(ctx context.Context, c Config) domain.Client {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: c.Endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.Key, c.Secret, "")),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load: %w", err))
	}

	s3Client := s3.NewFromConfig(cfg)

	return &client{
		client: s3Client,
		bucket: c.Bucket,
	}
}

func (s client) Upload(ctx context.Context, path string, reader io.Reader) error {
	var contentType string

	switch {
	case strings.Contains(path, ".svg"):
		contentType = "image/jpeg"
	case strings.Contains(path, ".mp4"):
		contentType = "video/mp4"
	case strings.Contains(path, ".jpg") || strings.Contains(path, ".jpeg"):
		contentType = "image/jpeg"
	case strings.Contains(path, ".png"):
		contentType = "image/png"
	case strings.Contains(path, ".webp"):
		contentType = "image/webp"
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	readSeeker := bytes.NewReader(buf.Bytes())

	_, err = s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(path),
			Body:        readSeeker,
			ContentType: aws.String(contentType),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to put: %w", err)
	}

	return nil
}
