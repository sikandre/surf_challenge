package action

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"surf_challenge/internal/action"
)

func Test_actionsHandler_GetNextActionProbability(t *testing.T) {
	type mocks struct {
		logger  *zap.SugaredLogger
		service *action.MockService
	}

	tests := []struct {
		name       string
		mock       func(m *mocks)
		nextInput  string
		wantStatus int
		assertBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Should return next action probability sorted successfully when valid next action is provided",
			nextInput: "action1",
			mock: func(m *mocks) {
				m.service.EXPECT().GetNextActionProbability(gomock.Any(), "action1").Return(
					map[string]string{
						"action3": "0.20", // not in order
						"action2": "0.70",
						"action4": "0.10",
					}, nil,
				)
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				wantBody := `{"action2":0.70,"action3":0.20,"action4":0.10}` + "\n"
				require.Equal(t, wantBody, recorder.Body.String())
			},
		},
		{
			name:      "Should return bad request when next action is not provided",
			nextInput: "",
			mock: func(m *mocks) {
				// No service call expected
			},
			wantStatus: http.StatusBadRequest,
			assertBody: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				wantBody := "next action parameter is required\n"
				require.Equal(t, wantBody, recorder.Body.String())
			},
		},
		{
			name:      "Should return internal server error when service returns an error",
			nextInput: "action1",
			mock: func(m *mocks) {
				m.service.EXPECT().GetNextActionProbability(gomock.Any(), "action1").Return(
					nil, assert.AnError,
				)
			},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				wantBody := "Internal server error\n"
				require.Equal(t, wantBody, recorder.Body.String())
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
					service: action.NewMockService(ctrl),
				}

				tt.mock(m)

				rctx := chi.NewRouteContext()

				u, _ := url.Parse("/actions/next-probability?next=" + tt.nextInput)

				req, err := http.NewRequestWithContext(
					context.WithValue(t.Context(), chi.RouteCtxKey, rctx),
					http.MethodGet,
					u.String(),
					nil,
				)
				require.NoError(t, err)

				recorder := httptest.NewRecorder()
				h := NewHandler(m.logger, m.service)
				h.GetNextActionProbability().ServeHTTP(recorder, req)

				require.Equal(t, tt.wantStatus, recorder.Code)
				tt.assertBody(t, recorder)
			},
		)
	}
}
