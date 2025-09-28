package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"surf_challenge/internal/api/apierror"
	"surf_challenge/internal/api/user/dto"
	"surf_challenge/internal/api/user/mapper"
	"surf_challenge/internal/user"
	"surf_challenge/internal/user/domain"
)

type Handler interface {
	GetUsers() http.HandlerFunc
}

type usersHandler struct {
	logger  *zap.SugaredLogger
	service user.Service
}

func NewHandler(sugar *zap.SugaredLogger, service user.Service) Handler {
	return &usersHandler{
		logger:  sugar,
		service: service,
	}
}

func (h *usersHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.handleGetUsers(r)
		if err != nil {
			h.logger.Errorw("failed to get users", "error", err)

			apiError := mapper.MapErrors(err)
			http.Error(w, apiError.Message, apiError.Code)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.logger.Errorw("failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (h *usersHandler) handleGetUsers(r *http.Request) (*dto.UsersResponse, error) {
	ctx := r.Context()

	userID, page, size, err := extractQueryParams(r)
	if err != nil {
		return nil, err
	}

	users, pagination, err := h.service.QueryUsers(
		ctx,
		domain.Query{
			ID:       userID,
			Page:     page,
			PageSize: size,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}

	usersResponse := mapper.MapUsersToDTO(users)
	paginationDTO := mapper.MapPaginationToDTO(pagination, page, size)

	resp := &dto.UsersResponse{
		Users:      usersResponse,
		Pagination: paginationDTO,
	}

	return resp, nil
}

func extractQueryParams(r *http.Request) (*int64, int, int, error) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		return nil, 0, 0, nil
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, 0, 0, apierror.NewAPIError("invalid userId parameter", http.StatusBadRequest)
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, 0, 0, apierror.NewAPIError("invalid page parameter", http.StatusBadRequest)
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return nil, 0, 0, apierror.NewAPIError("invalid pageSize parameter", http.StatusBadRequest)
	}

	return &userID, page, pageSize, nil
}
