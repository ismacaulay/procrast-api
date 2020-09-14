package models

import (
	"errors"

	"github.com/google/uuid"
)

type Item struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Created     *int64     `json:"created,omitempty"`
	Modified    *int64     `json:"modified,omitempty"`
}

func (l Item) Validate(required bool) error {
	if required {
		if l.Title == nil || l.Description == nil {
			return errors.New("")
		}
	}

	return nil
}
