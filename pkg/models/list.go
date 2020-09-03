package models

import (
	"errors"

	"github.com/google/uuid"
)

type List struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func (l List) Validate(required bool) error {
	if required {
		if l.Title == nil || l.Description == nil {
			return errors.New("")
		}
	}

	return nil
}

type Item struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
}
