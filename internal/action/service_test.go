package action

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"surf_challenge/internal/action/domain"
	"surf_challenge/internal/action/storage"
	"surf_challenge/internal/action/storage/entity"
)

func Test_service_GetActionByUserID(t *testing.T) {
	type mocks struct {
		logger *zap.SugaredLogger
		repo   *storage.MockRepository
	}

	tests := []struct {
		name    string
		userID  int64
		mock    func(m *mocks)
		want    []*domain.Action
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "should return actions successfully",
			userID: 1,
			mock: func(m *mocks) {
				m.repo.EXPECT().GetActionsByUserID(gomock.Any(), int64(1)).Return(
					[]*entity.Action{
						{
							ID:         1,
							Type:       "click",
							UserID:     1,
							TargetUser: 2,
							CreatedAt:  "2023-10-01T10:00:00Z",
						},
					}, nil,
				)
			},
			want: []*domain.Action{
				{
					ID:         1,
					Type:       "click",
					UserID:     1,
					TargetUser: 2,
					CreatedAt:  time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "should return error when repo fails",
			userID: 1,
			mock: func(m *mocks) {
				m.repo.EXPECT().GetActionsByUserID(gomock.Any(), int64(1)).Return(
					nil, assert.AnError,
				)
			},
			want:    nil,
			wantErr: assert.Error,
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

				s := &service{
					logger: m.logger,
					repo:   m.repo,
				}

				got, err := s.GetActionByUserID(t.Context(), tt.userID)

				assert.Equal(t, tt.want, got)
				tt.wantErr(t, err)
			},
		)
	}
}

func Test_service_GetNextActionProbability(t *testing.T) {
	type mocks struct {
		logger *zap.SugaredLogger
		repo   *storage.MockRepository
	}
	tests := []struct {
		name    string
		action  string
		mock    func(m *mocks)
		want    map[string]string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "should return next action probabilities successfully",
			action: "action_1",
			mock: func(m *mocks) {
				m.repo.EXPECT().GetAllActions(gomock.Any()).Return(
					[]*entity.Action{
						{ID: 1, Type: "action_1", UserID: 1, TargetUser: 2, CreatedAt: "2023-10-01T10:00:00Z"},
						{ID: 2, Type: "action_2", UserID: 1, TargetUser: 3, CreatedAt: "2023-10-01T11:00:00Z"},
						{ID: 3, Type: "action_1", UserID: 2, TargetUser: 3, CreatedAt: "2023-10-01T12:00:00Z"},
						{ID: 4, Type: "action_3", UserID: 2, TargetUser: 4, CreatedAt: "2023-10-01T13:00:00Z"},
						{ID: 5, Type: "action_1", UserID: 3, TargetUser: 4, CreatedAt: "2023-10-01T14:00:00Z"},
						{ID: 6, Type: "action_2", UserID: 3, TargetUser: 5, CreatedAt: "2023-10-01T15:00:00Z"},
					}, nil,
				)
			},
			want: map[string]string{
				"action_2": "0.67",
				"action_3": "0.33",
			},
			wantErr: assert.NoError,
		},
		{
			name:   "should return empty map when no next actions found",
			action: "action_4",
			mock: func(m *mocks) {
				m.repo.EXPECT().GetAllActions(gomock.Any()).Return(
					[]*entity.Action{
						{ID: 1, Type: "action_1", UserID: 1, TargetUser: 2, CreatedAt: "2023-10-01T10:00:00Z"},
						{ID: 2, Type: "action_2", UserID: 1, TargetUser: 3, CreatedAt: "2023-10-01T11:00:00Z"},
					}, nil,
				)
			},
			want:    map[string]string{},
			wantErr: assert.NoError,
		},
		{
			name:   "should return error when repo fails",
			action: "action_1",
			mock: func(m *mocks) {
				m.repo.EXPECT().GetAllActions(gomock.Any()).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:   "should return empty map when no actions in repo",
			action: "action_1",
			mock: func(m *mocks) {
				m.repo.EXPECT().GetAllActions(gomock.Any()).Return([]*entity.Action{}, nil)
			},
			want:    map[string]string{},
			wantErr: assert.NoError,
		},
		{
			name:   "should return empty map when only one action in repo",
			action: "action_1",
			mock: func(m *mocks) {
				m.repo.EXPECT().GetAllActions(gomock.Any()).Return(
					[]*entity.Action{
						{ID: 1, Type: "action_1", UserID: 1, TargetUser: 2, CreatedAt: "2023-10-01T10:00:00Z"},
					}, nil,
				)
			},
			want:    map[string]string{},
			wantErr: assert.NoError,
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

				s := &service{
					logger: m.logger,
					repo:   m.repo,
				}

				got, err := s.GetNextActionProbability(t.Context(), tt.action)

				assert.Equal(t, tt.want, got)
				tt.wantErr(t, err)
			},
		)
	}
}
