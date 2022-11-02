package domain

import "github.com/vediagames/zeroerror"

type User struct {
	Email string
	Name  string
}

func (u User) Validate() error {
	var err zeroerror.Error

	if u.Email == "" {
		err.Add(ErrEmptyEmail)
	}

	if u.Name == "" {
		err.Add(ErrEmptyName)
	}

	return err.Err()
}
