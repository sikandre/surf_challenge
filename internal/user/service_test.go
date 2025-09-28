package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"surf_challenge/internal/action"
	actiondomain "surf_challenge/internal/action/domain"
	"surf_challenge/internal/user/domain"
	"surf_challenge/internal/user/storage"
	"surf_challenge/internal/user/storage/entity"
)

func Test_userService_QueryUsers(t *testing.T) {
	type mocks struct {
		logger *zap.SugaredLogger
		repo   *storage.MockRepository
	}

	tests := []struct {
		name        string
		mock        func(m *mocks)
		query       domain.Query
		want        []*domain.User
		wantResults *domain.Results
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "should return users successfully",
			query: domain.Query{
				ID:       nil,
				Page:     1,
				PageSize: 10,
			},
			mock: func(m *mocks) {
				m.repo.EXPECT().QueryUsers(gomock.Any(), nil, 1, 10).Return(
					[]*entity.User{
						{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: "2023-10-01T10:00:00Z",
						},
						{
							ID:        2,
							Name:      "Jane Smith",
							CreatedAt: "2023-10-02T11:00:00Z",
						},
					}, 2, nil,
				)
			},
			want: []*domain.User{
				{
					ID:        1,
					Name:      "John Doe",
					CreatedAt: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					Name:      "Jane Smith",
					CreatedAt: time.Date(2023, 10, 2, 11, 0, 0, 0, time.UTC),
				},
			},
			wantResults: &domain.Results{
				TotalItems: 2,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return error when repository fails",
			query: domain.Query{
				ID:       nil,
				Page:     1,
				PageSize: 10,
			},
			mock: func(m *mocks) {
				m.repo.EXPECT().QueryUsers(gomock.Any(), nil, 1, 10).Return(
					nil, 0, assert.AnError,
				)
			},
			want:        nil,
			wantResults: nil,
			wantErr:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				m := &mocks{
					logger: zap.NewNop().Sugar(),
					repo:   storage.NewMockRepository(ctrl),
				}

				tt.mock(m)

				s := &userService{
					logger: m.logger,
					repo:   m.repo,
				}
				got, results, err := s.QueryUsers(t.Context(), tt.query)

				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantResults, results)
				tt.wantErr(t, err)
			},
		)
	}
}

func Test_userService_GetUserActionCount(t *testing.T) {
	type mocks struct {
		logger        *zap.SugaredLogger
		repo          *storage.MockRepository
		actionService *action.MockService
	}

	tests := []struct {
		name    string
		userID  int64
		mock    func(m *mocks)
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "should return action count successfully",
			userID: 1,
			mock: func(m *mocks) {
				m.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(
						&entity.User{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: "2023-10-01T10:00:00Z",
						}, nil,
					)

				m.actionService.EXPECT().
					GetActionByUserID(gomock.Any(), int64(1)).
					Return(
						[]*actiondomain.Action{
							{
								ID:         1,
								Type:       "click",
								UserID:     1,
								TargetUser: 0,
								CreatedAt:  "2023-10-05T12:00:00Z",
							},
						}, nil,
					)
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name:   "should return not found error when user does not exist",
			userID: 99,
			mock: func(m *mocks) {
				m.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(99)).
					Return(nil, storage.ErrUserNotFound)
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNotFound)
			},
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			mock: func(m *mocks) {
				m.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(nil, assert.AnError)
			},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name:   "should return error when action service fails",
			userID: 1,
			mock: func(m *mocks) {
				m.repo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(
						&entity.User{
							ID:        1,
							Name:      "John Doe",
							CreatedAt: "2023-10-01T10:00:00Z",
						}, nil,
					)

				m.actionService.EXPECT().
					GetActionByUserID(gomock.Any(), int64(1)).
					Return(nil, assert.AnError)
			},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				m := &mocks{
					logger:        zap.NewNop().Sugar(),
					repo:          storage.NewMockRepository(ctrl),
					actionService: action.NewMockService(ctrl),
				}
				tt.mock(m)

				s := &userService{
					logger:        m.logger,
					repo:          m.repo,
					actionService: m.actionService,
				}
				got, err := s.GetUserActionCount(t.Context(), tt.userID)

				assert.Equal(t, tt.want, got)
				tt.wantErr(t, err)
			},
		)
	}
}
