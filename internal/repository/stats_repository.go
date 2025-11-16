package repository

import (
	"context"
	"fmt"

	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

type StatsRepository interface {
	GetReviewStats(ctx context.Context) ([]api.UserReviewStat, error)
}

func (r *PostgresRepository) GetReviewStats(ctx context.Context) ([]api.UserReviewStat, error) {
	var stats []api.UserReviewStat

	rawQuery := `
		SELECT reviewer_id as user_id, COUNT(*) as review_count
		FROM (
			SELECT unnest(string_to_array(assigned_reviewers, ',')) as reviewer_id
			FROM pull_requests
			WHERE assigned_reviewers != ''
		) as reviewers
		GROUP BY reviewer_id
		ORDER BY review_count DESC;
	`

	if err := r.DB.WithContext(ctx).Raw(rawQuery).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении статистики по ревью: %w", err)
	}

	return stats, nil
}
