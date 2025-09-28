package domain

type Action struct {
	ID         int
	Type       string
	UserID     int
	TargetUser int
	CreatedAt  string
}
