package domain

import "time"

type Session struct {
	ID         string
	CreatedAt  time.Time
	InsertedAt time.Time
}
