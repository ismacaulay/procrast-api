package models

import (
	"github.com/google/uuid"
)

type List struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     int64     `json:"created"`
	Modified    int64     `json:"modified"`
}

type Item struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       uint8     `json:"state"`
	Created     int64     `json:"created"`
	Modified    int64     `json:"modified"`
	ListUUID    uuid.UUID `json:"list_uuid"`
}

type History struct {
	UUID      uuid.UUID `json:"uuid"`
	Command   string    `json:"command"`
	State     []byte    `json:"state"`
	Timestamp int64     `json:"timestamp"`
	Created   int64     `json:"created"`
}

type User struct {
	UUID     uuid.UUID
	Email    string
	PassHash []byte
	Created  int64
	Modified int64
}
