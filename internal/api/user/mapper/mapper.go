package mapper

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"surf_challenge/internal/api/apierror"
	"surf_challenge/internal/api/user/dto"
	"surf_challenge/internal/user"
	"surf_challenge/internal/user/domain"
)

func MapErrors(err error) *apierror.APIError {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	switch {
	case errors.Is(err, user.ErrNotFound):
		return apierror.NewAPIError("Resource not found", http.StatusNotFound)
	default:
		return apierror.NewAPIError("Internal server error", http.StatusInternalServerError)
	}
}

func MapUsersToDTO(users []*domain.User) []dto.User {
	userDTOs := make([]dto.User, len(users))
	for i, u := range users {
		userDTOs[i] = MapUserToDTO(u)
	}

	return userDTOs
}

func MapPaginationToDTO(pagination *domain.Results, page int, size int) dto.Pagination {
	totalPages := (pagination.TotalItems + size - 1) / size

	return dto.Pagination{
		TotalItems: pagination.TotalItems,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   size,
	}
}

func MapUserToDTO(u *domain.User) dto.User {
	id := strconv.Itoa(int(u.ID))

	return dto.User{
		ID:        id,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
