package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/model"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"gorm.io/gorm"
)

type PullRequestRepository interface {
	SavePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error)
	GetPullRequest(ctx context.Context, prId string) (api.PullRequest, error)
}

func (r *PostgresRepository) SavePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error) {
	var existingPR model.PullRequest
	if r.DB.WithContext(ctx).Where("pull_request_id = ?", pr.PullRequestId).First(&existingPR).RowsAffected > 0 {
		return api.PullRequest{}, fmt.Errorf("%w: Пул реквест с ID %s уже существует", errWrappers.ErrPrExists, pr.PullRequestId)
	}

	pr.CreatedAt = func() *time.Time { t := time.Now(); return &t }()
	pr.Status = api.PullRequestStatusOPEN

	prModel := model.FromAPIPullRequest(pr)
	if err := r.DB.WithContext(ctx).Create(&prModel).Error; err != nil {
		return api.PullRequest{}, err
	}
	return prModel.ToAPIPullRequest(), nil
}

func (r *PostgresRepository) UpdatePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error) {
	pullRequestModel := model.FromAPIPullRequest(pr)

	if err := r.DB.WithContext(ctx).Where("pull_request_id = ?", pr.PullRequestId).Updates(pullRequestModel).Error; err != nil {
		return api.PullRequest{}, err
	}
	return r.GetPullRequest(ctx, pr.PullRequestId)
}

func (r *PostgresRepository) GetPullRequest(ctx context.Context, prId string) (api.PullRequest, error) {
	var pullRequestModel model.PullRequest
	if err := r.DB.WithContext(ctx).Where("pull_request_id = ?", prId).First(&pullRequestModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.PullRequest{}, fmt.Errorf("%w: Пул реквест с ID %s не существует", errWrappers.ErrNotFound, prId)
		}
		return api.PullRequest{}, err
	}
	return pullRequestModel.ToAPIPullRequest(), nil
}
