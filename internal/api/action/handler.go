package action

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"surf_challenge/internal/action"
	"surf_challenge/internal/api/action/dto"
	"surf_challenge/internal/api/action/mapper"
	"surf_challenge/internal/api/apierror"
)

type Handler interface {
	GetNextActionProbability() http.HandlerFunc
}

type actionsHandler struct {
	logger  *zap.SugaredLogger
	service action.Service
}

func NewHandler(sugar *zap.SugaredLogger, service action.Service) Handler {
	return &actionsHandler{
		logger:  sugar,
		service: service,
	}
}

func (a actionsHandler) GetNextActionProbability() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := a.handleGetNextActionProbability(r)
		if err != nil {
			a.logger.Errorw("failed to get next action probability", "error", err)

			http.Error(w, "Failed to get next action probability", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			a.logger.Errorw("failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (a actionsHandler) handleGetNextActionProbability(r *http.Request) (*dto.NextActionProbability, error) {
	ctx := r.Context()

	nextAction := r.URL.Query().Get("next")
	if nextAction == "" {
		return nil, apierror.NewAPIError("next action parameter is required", http.StatusBadRequest)
	}

	probability, err := a.service.GetNextActionProbability(ctx, nextAction)
	if err != nil {
		return nil, err
	}

	probabilityDTO, err := mapper.MapProbabilityToDTO(probability)
	if err != nil {
		return nil, err
	}

	return probabilityDTO, nil
}
