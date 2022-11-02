package domain

import (
	"context"
	"fmt"

	"github.com/vediagames/zeroerror"
)

type EmailClient interface {
	Email(context.Context, EmailRequest) error
}

type EmailRequest struct {
	To      User
	From    User
	Name    string
	Subject string
	Body    string
}

func (e EmailRequest) Validate() error {
	var err zeroerror.Error

	if ve := e.To.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid to: %w", ve))
	}

	if ve := e.From.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid from: %w", ve))
	}

	if e.Subject == "" {
		err.Add(ErrEmptySubject)
	}

	if e.Body == "" {
		err.Add(ErrEmptyBody)
	}

	return err.Err()
}
