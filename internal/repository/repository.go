package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/model"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"gorm.io/gorm"
)

type Repository interface {
	SaveTeam(ctx context.Context, team api.Team) (api.Team, error)
	GetTeam(ctx context.Context, teamName string) (api.Team, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]model.User, error)
	GetUser(ctx context.Context, userId string) (api.User, error)
	SetUserIsActive(ctx context.Context, userId string, isActive bool) (api.User, error)

	SavePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr api.PullRequest) (api.PullRequest, error)
	GetPullRequest(ctx context.Context, prId string) (api.PullRequest, error)

	FindUserPullRequests(ctx context.Context, userId string) ([]api.PullRequestShort, error)

	FindActiveCandidates(ctx context.Context, teamName string, excludeIds []string) ([]string, error)
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) SaveTeam(ctx context.Context, team api.Team) (api.Team, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		teamModel := model.Team{TeamName: team.TeamName}
		if result := tx.Where("team_name = ?", team.TeamName).First(&teamModel); result.RowsAffected > 0 {
			return fmt.Errorf("%w: Команда с именем %s уже существует", errWrappers.ErrTeamExists, team.TeamName)
		}
		if err := tx.Create(&teamModel).Error; err != nil {
			return err
		}

		for _, member := range team.Members {
			userModel := model.User{
				UserId:   member.UserId,
				Username: member.Username,
				TeamName: team.TeamName,
				IsActive: member.IsActive,
			}

			if err := tx.Create(&userModel).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return api.Team{}, err
	}
	return team, nil
}

func (r *PostgresRepository) GetTeam(ctx context.Context, teamName string) (api.Team, error) {
	var teamModel model.Team
	if err := r.DB.WithContext(ctx).Where("team_name = ?", teamName).First(&teamModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.Team{}, fmt.Errorf("%w: Команда с именем %s не найдена", errWrappers.ErrNotFound, teamName)
		}
		return api.Team{}, err
	}

	members, _ := r.GetTeamMembers(ctx, teamName)
	apiMembers := make([]api.TeamMember, len(members))

	for i, member := range members {
		apiMembers[i] = member.ToAPITeamMember()
	}

	return api.Team{TeamName: teamName, Members: apiMembers}, nil
}

func (r *PostgresRepository) GetTeamMembers(ctx context.Context, teamName string) ([]model.User, error) {
	var members []model.User
	err := r.DB.WithContext(ctx).Where("team_name = ?", teamName).Find(&members).Error
	return members, err
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

	if err := r.DB.WithContext(ctx).Where("pull_request_id - ?", pr.PullRequestId).Updates(pullRequestModel).Error; err != nil {
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

func (r *PostgresRepository) FindUserPullRequests(ctx context.Context, userId string) ([]api.PullRequestShort, error) {
	var pullRequestModels []model.PullRequest

	searchPattern := "%" + userId + "%"
	if err := r.DB.WithContext(ctx).Where("assigned_reviewers LIKE ?", searchPattern).Find(pullRequestModels).Error; err != nil {
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

func (r *PostgresRepository) FindActiveCandidates(ctx context.Context, teamName string, excludeIds []string) ([]string, error) {
	var userModels []model.User

	query := r.DB.WithContext(ctx).Where("team_name = ? AND is_active = ?", teamName, true)

	if len(excludeIds) > 0 {
		query = query.Not("user_id", excludeIds)
	}

	if err := query.Find(&userModels).Error; err != nil {
		return nil, err
	}

	candidates := make([]string, len(userModels))
	for i, user := range userModels {
		candidates[i] = user.UserId
	}
	return candidates, nil
}

func (r *PostgresRepository) ChooseRandomCandidates(candidates []string, count int) []string {
	if len(candidates) <= count {
		return candidates
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := random.Perm(len(candidates))
	chosen := make([]string, 0, count)
	for i := 0; i < count; i++ {
		chosen = append(chosen, candidates[perm[i]])
	}
	return chosen
}
