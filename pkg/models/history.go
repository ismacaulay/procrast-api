package models

import "github.com/google/uuid"

type History struct {
	Id      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	State   []byte    `json:"state"`
	Created int64     `json:"created"`
}
