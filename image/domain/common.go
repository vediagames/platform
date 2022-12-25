package domain

import (
	"github.com/vediagames/zeroerror"
)

type Image struct {
	ID int
}

func (im Image) Validate() error {
	var err zeroerror.Error

	if im.ID < 0 {
		err.Add(ErrInvalidID)
	}

	return err.Err()
}
