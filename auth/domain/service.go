package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/vediagames/zeroerror"
)

type Service interface {
	Authenticate(context.Context, AuthenticateRequest) (AuthenticateResponse, error)
	ToContext(context.Context, User) context.Context
	FromContext(context.Context) (User, error)
}

type AuthenticateRequest struct {
	Cookies string
}

func (r AuthenticateRequest) Validate() error {
	var err zeroerror.Error

	if r.Cookies == "" {
		err.Add(fmt.Errorf("empty cookies"))
	}

	return err.Err()
}

type AuthenticateResponse struct {
	User User
}

func (r AuthenticateResponse) Validate() error {
	var err zeroerror.Error

	if ve := r.User.Validate(); ve != nil {
		err.Add(fmt.Errorf("invalid user: %w", err))
	}

	return err.Err()
}

type User struct {
	ID        string
	SessionID string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Validate() error {
	var err zeroerror.Error

	if u.ID == "" {
		err.Add(fmt.Errorf("empty ID"))
	}

	if u.SessionID == "" {
		err.Add(fmt.Errorf("empty session ID"))
	}

	if u.Username == "" {
		err.Add(fmt.Errorf("empty username"))
	}

	if u.Email == "" {
		err.Add(fmt.Errorf("empty email"))
	}

	if u.CreatedAt.IsZero() {
		err.Add(fmt.Errorf("zero created at"))
	}

	if u.UpdatedAt.IsZero() {
		err.Add(fmt.Errorf("zero updated at"))
	}

	return err.Err()
}
