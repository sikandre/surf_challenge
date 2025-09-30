package mapper

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"surf_challenge/internal/api/action/dto"
	"surf_challenge/internal/api/apierror"
)

func MapProbabilityToDTO(probability map[string]string) (*dto.NextActionProbability, error) {
	keys := make([]string, 0, len(probability))

	floatMapData := make(map[string]float64, len(probability))
	for k, v := range probability {
		keys = append(keys, k)

		prob, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse probability value %s for key %s: %w", v, k, err)
		}

		floatMapData[k] = prob
	}

	sort.Slice(
		keys, func(i, j int) bool {
			return floatMapData[keys[i]] > floatMapData[keys[j]]
		},
	)

	return &dto.NextActionProbability{
		Data: floatMapData,
		Keys: keys,
	}, nil
}

func MapErrors(err error) *apierror.APIError {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	return apierror.NewAPIError("Internal server error", http.StatusInternalServerError)
}
