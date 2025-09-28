package domain

import "time"

type Results struct {
	TotalItems int
}

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type Query struct {
	ID       *int64
	Page     int
	PageSize int
}
