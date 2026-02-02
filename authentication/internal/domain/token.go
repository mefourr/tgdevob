package domain

import "time"

type Token struct {
	IamToken  string
	ExpiresAt time.Time
}
