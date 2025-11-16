package handler

import (
	"context"

	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

func (s *Server) GetStatsReviews(ctx context.Context, request api.GetStatsReviewsRequestObject) (api.GetStatsReviewsResponseObject, error) {
	stats, err := s.Repository.GetReviewStats(ctx)
	if err != nil {
		return nil, err
	}

	return api.GetStatsReviews200JSONResponse{Stats: stats}, nil
}
