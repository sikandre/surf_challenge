package storage

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"sync"

	"surf_challenge/internal/action/storage/entity"
)

var ErrActionsNotFound = errors.New("actions not found")

var (
	actionsOnce  sync.Once
	actionsCache []*entity.Action
	errActions   error
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=storage
type Repository interface {
	GetActionsByUserID(ctx context.Context, userID int64) ([]*entity.Action, error)
	GetAllActions(ctx context.Context) ([]*entity.Action, error)
}

type actionRepository struct{}

func NewRepository() Repository {
	return &actionRepository{}
}

func (ar *actionRepository) GetAllActions(_ context.Context) ([]*entity.Action, error) {
	actions, err := loadFileWithActions()
	if err != nil {
		return nil, err
	}

	if len(actions) == 0 {
		return nil, ErrActionsNotFound
	}

	return actions, nil
}

func (ar *actionRepository) GetActionsByUserID(_ context.Context, userID int64) ([]*entity.Action, error) {
	actions, err := loadFileWithActions()
	if err != nil {
		return nil, err
	}

	var filteredActions []*entity.Action

	for _, action := range actions {
		if int64(action.UserID) == userID {
			filteredActions = append(filteredActions, action)
		}
	}

	if len(filteredActions) == 0 {
		return nil, ErrActionsNotFound
	}

	return filteredActions, nil
}

//go:embed db/actions.json
var actionsFile []byte

func loadFileWithActions() ([]*entity.Action, error) {
	actionsOnce.Do(
		func() {
			actionsCache, errActions = parseActionsFile()
		},
	)

	return actionsCache, errActions
}

func parseActionsFile() ([]*entity.Action, error) {
	var actions []*entity.Action

	err := json.Unmarshal(actionsFile, &actions)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
