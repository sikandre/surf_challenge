package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"surf_challenge/internal/api/user/dto"
	"surf_challenge/internal/converter"
	"surf_challenge/internal/user"
	"surf_challenge/internal/user/domain"
)

func Test_usersHandler_GetUsers(t *testing.T) {
	type mocks struct {
		logger  *zap.SugaredLogger
		service *user.MockService
	}

	tests := []struct {
		name       string
		mock       func(m *mocks)
		query      map[string]string
		wantStatus int
		assertBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "When no query params is provided, should return all users with default pagination",
			mock: func(m *mocks) {
				m.service.EXPECT().QueryUsers(
					gomock.Any(),
					domain.Query{
						ID:       nil,
						Page:     1,
						PageSize: 10,
					},
				).Return(
					[]*domain.User{
						{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					&domain.Results{TotalItems: 1},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				want := &dto.UsersResponse{
					Users: []dto.User{
						{
							ID:        "1",
							Name:      "John Doe",
							CreatedAt: "2023-01-01T00:00:00Z",
						},
					},
					Pagination: dto.Pagination{
						TotalItems: 1,
						Page:       1,
						PageSize:   10,
						TotalPages: 1,
					},
				}

				expected, err := json.Marshal(want)
				require.NoError(t, err)
				assert.JSONEq(t, string(expected), r.Body.String())
			},
		},
		{
			name: "When user ID is provided, should return the user with that ID",
			query: map[string]string{
				"userId": "1",
			},
			mock: func(m *mocks) {
				m.service.EXPECT().QueryUsers(
					gomock.Any(),
					domain.Query{
						ID:       converter.ToPtr(int64(1)),
						Page:     1,
						PageSize: 10,
					},
				).Return(
					[]*domain.User{
						{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					&domain.Results{TotalItems: 1},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				want := &dto.UsersResponse{
					Users: []dto.User{
						{
							ID:        "1",
							Name:      "John Doe",
							CreatedAt: "2023-01-01T00:00:00Z",
						},
					},
					Pagination: dto.Pagination{
						TotalItems: 1,
						Page:       1,
						PageSize:   10,
						TotalPages: 1,
					},
				}

				expected, err := json.Marshal(want)
				require.NoError(t, err)
				assert.JSONEq(t, string(expected), r.Body.String())
			},
		},
		{
			name: "When user ID not provided and page and pageSize provided, should return users and pagination applied",
			query: map[string]string{
				"page":     "4",
				"pageSize": "1",
			},
			mock: func(m *mocks) {
				m.service.EXPECT().QueryUsers(
					gomock.Any(),
					domain.Query{
						ID:       nil,
						Page:     4,
						PageSize: 1,
					},
				).Return(
					[]*domain.User{
						{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					&domain.Results{TotalItems: 5},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				want := &dto.UsersResponse{
					Users: []dto.User{
						{
							ID:        "1",
							Name:      "John Doe",
							CreatedAt: "2023-01-01T00:00:00Z",
						},
					},
					Pagination: dto.Pagination{
						TotalItems: 5,
						Page:       4,
						PageSize:   1,
						TotalPages: 5,
					},
				}

				expected, err := json.Marshal(want)
				require.NoError(t, err)
				assert.JSONEq(t, string(expected), r.Body.String())
			},
		},
		{
			name: "When user ID is provided but not found, should return empty list",
			query: map[string]string{
				"userId": "1",
			},
			mock: func(m *mocks) {
				m.service.EXPECT().QueryUsers(
					gomock.Any(),
					domain.Query{
						ID:       converter.ToPtr(int64(1)),
						Page:     1,
						PageSize: 10,
					},
				).Return(
					nil,
					&domain.Results{TotalItems: 0},
					nil,
				)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				want := &dto.UsersResponse{
					Users: []dto.User{},
					Pagination: dto.Pagination{
						TotalItems: 0,
						Page:       1,
						PageSize:   10,
						TotalPages: 0,
					},
				}

				expected, err := json.Marshal(want)
				require.NoError(t, err)
				assert.JSONEq(t, string(expected), r.Body.String())
			},
		},
		{
			name:  "When service returns an error, should return internal server error",
			query: map[string]string{},
			mock: func(m *mocks) {
				m.service.EXPECT().QueryUsers(
					gomock.Any(),
					domain.Query{
						ID:       nil,
						Page:     1,
						PageSize: 10,
					},
				).Return(
					nil,
					nil,
					assert.AnError,
				)
			},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				assert.Contains(t, r.Body.String(), "Internal server error")
			},
		},
		{
			name: "When user ID is not an integer, should return bad request",
			query: map[string]string{
				"userId": "abc",
			},
			mock:       func(m *mocks) {},
			wantStatus: http.StatusBadRequest,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				assert.Contains(t, r.Body.String(), "invalid userId parameter")
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				m := &mocks{
					logger:  zap.NewNop().Sugar(),
					service: user.NewMockService(ctrl),
				}

				tt.mock(m)

				rctx := chi.NewRouteContext()

				u, _ := url.Parse("/api/v1/users")

				vals := u.Query()
				for k, v := range tt.query {
					vals.Set(k, v)
				}

				u.RawQuery = vals.Encode()

				req, err := http.NewRequestWithContext(
					context.WithValue(t.Context(), chi.RouteCtxKey, rctx),
					http.MethodGet,
					u.String(),
					nil,
				)
				require.NoError(t, err)

				recorder := httptest.NewRecorder()
				h := NewHandler(m.logger, m.service)
				h.GetUsers().ServeHTTP(recorder, req)

				tt.assertBody(t, recorder)
			},
		)
	}
}

func Test_usersHandler_GetUserActionCount(t *testing.T) {
	type mocks struct {
		logger  *zap.SugaredLogger
		service *user.MockService
	}

	tests := []struct {
		name       string
		userID     string
		mock       func(m *mocks)
		wantStatus int
		assertBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "When user ID is valid, should return action count",
			userID: "1",
			mock: func(m *mocks) {
				m.service.EXPECT().GetUserActionCount(
					gomock.Any(),
					int64(1),
				).Return(5, nil)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				want := dto.ActionsCount{
					Count: 5,
				}

				expected, err := json.Marshal(want)
				require.NoError(t, err)
				assert.JSONEq(t, string(expected), r.Body.String())
			},
		},
		{
			name:       "When user ID is not an integer, should return bad request",
			userID:     "abc",
			mock:       func(m *mocks) {},
			wantStatus: http.StatusBadRequest,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				assert.Contains(t, r.Body.String(), "invalid userId parameter")
			},
		},
		{
			name:   "When service returns an error, should return internal server error",
			userID: "1",
			mock: func(m *mocks) {
				m.service.EXPECT().GetUserActionCount(
					gomock.Any(),
					int64(1),
				).Return(0, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				assert.Contains(t, r.Body.String(), "Internal server error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				m := &mocks{
					logger:  zap.NewNop().Sugar(),
					service: user.NewMockService(ctrl),
				}

				tt.mock(m)

				rctx := chi.NewRouteContext()

				u, _ := url.Parse("/api/v1/users/" + tt.userID + "/actions/count")
				rctx.URLParams.Add("userId", tt.userID)

				req, err := http.NewRequestWithContext(
					context.WithValue(t.Context(), chi.RouteCtxKey, rctx),
					http.MethodGet,
					u.String(),
					nil,
				)
				require.NoError(t, err)

				recorder := httptest.NewRecorder()
				h := NewHandler(m.logger, m.service)
				h.GetUserActionCount().ServeHTTP(recorder, req)

				tt.assertBody(t, recorder)
			},
		)
	}
}
