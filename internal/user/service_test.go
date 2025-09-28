package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

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
