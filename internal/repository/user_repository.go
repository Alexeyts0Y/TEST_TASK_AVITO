package repository

import (
	"context"
	"errors"
	"fmt"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/model"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(ctx context.Context, userId string) (api.User, error)
	SetUserIsActive(ctx context.Context, userId string, isActive bool) (api.User, error)
	FindUserPullRequests(ctx context.Context, userId string) ([]api.PullRequestShort, error)
}

func (r *PostgresRepository) GetUser(ctx context.Context, userId string) (api.User, error) {
	var userModel model.User
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userId).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.User{}, fmt.Errorf("%w: Пользователь с ID %s не найден", errWrappers.ErrNotFound, userId)
		}
		return api.User{}, err
	}

	return api.User{
		UserId:   userModel.UserId,
		Username: userModel.Username,
		TeamName: userModel.TeamName,
		IsActive: userModel.IsActive}, nil
}

func (r *PostgresRepository) SetUserIsActive(ctx context.Context, userId string, isActive bool) (api.User, error) {
	user, err := r.GetUser(ctx, userId)
	if err != nil {
		return api.User{}, err
	}

	if err := r.DB.WithContext(ctx).Model(&model.User{}).Where("user_id = ?", userId).Update("is_active", isActive).Error; err != nil {
		return api.User{}, err
	}

	user.IsActive = isActive
	return user, nil
}

func (r *PostgresRepository) FindUserPullRequests(ctx context.Context, userId string) ([]api.PullRequestShort, error) {
	var pullRequestModels []model.PullRequest

	searchPattern := "%" + userId + "%"
	if err := r.DB.WithContext(ctx).Where("assigned_reviewers LIKE ?", searchPattern).Find(&pullRequestModels).Error; err != nil {
		return nil, err
	}

	shortPullRequests := make([]api.PullRequestShort, len(pullRequestModels))
	for i, model := range pullRequestModels {
		shortPullRequests[i] = api.PullRequestShort{
			AuthorId:        model.AuthorId,
			PullRequestId:   model.PullRequestId,
			PullRequestName: model.PullRequestName,
			Status:          api.PullRequestShortStatus(model.Status),
		}
	}
	return shortPullRequests, nil
}
