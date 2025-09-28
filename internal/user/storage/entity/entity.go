package entity

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"` // in ISO 8601 format (e.g., "2022-04-14T11:12:22.758Z") RFC3339
}
