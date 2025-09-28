package action

import (
	"testing"

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
					CreatedAt:  "2023-10-01T10:00:00Z",
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
