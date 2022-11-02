package sendinblue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/vediagames/vediagames.com/notification/domain"
)

type emailRequest struct {
	Subject     string `json:"subject"`
	Sender      sender `json:"sender"`
	To          []to   `json:"to"`
	HTMLContent string `json:"htmlContent"`
}

type sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type to struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type service struct {
	client http.Client
	token  string
}

func New(client http.Client, token string) domain.EmailClient {
	return &service{
		client: client,
		token:  token,
	}
}

func (s service) Email(ctx context.Context, req domain.EmailRequest) error {
	body := emailRequest{
		Subject: req.Subject,
		Sender:  sender(req.From),
		To: []to{
			to(req.To),
		},
		HTMLContent: req.Body,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal email body: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, "https://api.sendinblue.com/v3/smtp/email", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api-key", s.token)
	httpReq.Header.Set("accept", "application/json")

	httpRes, err := s.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer httpRes.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	zerolog.
		Ctx(ctx).
		Info().
		Str("body", buf.String()).
		Int("status", httpRes.StatusCode).
		Str("component", "email").
		Str("provider", "sendinblue").
		Send()

	if httpRes.StatusCode >= 400 {
		return fmt.Errorf("http request failed with code: %v", httpRes.StatusCode)
	}

	return nil
}
