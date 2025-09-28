package entity

type Action struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	UserID     int    `json:"userId"`
	TargetUser int    `json:"targetUser"`
	CreatedAt  string `json:"createdAt"`
}
