package dto

type Pagination struct {
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"` // in ISO 8601 format (e.g., "2022-04-14T11:12:22.758Z") RFC3339
}

type UsersResponse struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pagination"`
}
