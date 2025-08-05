package dao

import "time"

type BatchToken struct {
	Token     string
	ExpiresAt time.Time
}
