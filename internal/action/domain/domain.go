package domain

import "time"

type Action struct {
	ID         int
	Type       string
	UserID     int
	TargetUser int
	CreatedAt  time.Time
}
