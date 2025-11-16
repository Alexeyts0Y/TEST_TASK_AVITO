package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/model"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/utils"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"gorm.io/gorm"
)

type TeamRepository interface {
	SaveTeam(ctx context.Context, team api.Team) (api.Team, error)
	GetTeam(ctx context.Context, teamName string) (api.Team, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]model.User, error)
	FindActiveCandidates(ctx context.Context, teamName string, excludeIds []string) ([]string, error)
	DeactivateTeamMembers(ctx context.Context, teamName string) (int64, error)
	ReassignPRsForTeam(ctx context.Context, teamName string) (api.ReassignmentSummary, error)
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

func (r *PostgresRepository) DeactivateTeamMembers(ctx context.Context, teamName string) (int64, error) {
	result := r.DB.WithContext(ctx).Model(&model.User{}).Where("team_name = ?", teamName).Update("is_active", false)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%w: команда с именем %s не найдена", errWrappers.ErrNotFound, teamName)
	} else if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *PostgresRepository) ReassignPRsForTeam(ctx context.Context, teamName string) (api.ReassignmentSummary, error) {
	var pullRequests []model.PullRequest
	var summary api.ReassignmentSummary
	summary.TeamName = teamName

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rawQuery := `
        SELECT pr.* 
        FROM pull_requests pr
        WHERE pr.status = 'OPEN'
          AND EXISTS (
            SELECT 1 
            FROM users u 
            WHERE u.team_name = ?
              AND u.user_id = ANY(string_to_array(pr.assigned_reviewers, ','))
          )`

		err := tx.Raw(rawQuery, teamName).Scan(&pullRequests).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: не найдены PR для команды %s", errWrappers.ErrNotFound, teamName)
		} else if err != nil {
			return fmt.Errorf("%w: ошибка при выборке PR для команды %s", err, teamName)
		}

		count := 0
		for _, pr := range pullRequests {
			excludeIds := strings.Split(pr.AssignedReviewers, ",")
			candidates, err := r.FindActiveCandidates(ctx, teamName, excludeIds)
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: не найдены PR для команды %s", errWrappers.ErrNotFound, teamName)
			} else if err != nil {
				return fmt.Errorf("%w: ошибка при поиске актвных кандидатов команды %s", err, teamName)
			}

			if len(candidates) == 0 {
				continue
			}

			currentReviewersCount := len(strings.Split(pr.AssignedReviewers, ","))
			chosen := utils.ChooseRandomCandidates(candidates, currentReviewersCount)
			if err := tx.Where("pull_request_id = ?", pr.PullRequestId).Model(&model.PullRequest{}).Update("assigned_reviewers", strings.Join(chosen, ",")).Error; err != nil {
				return err
			}
			count++
		}
		summary.ReassignedPrsCount = count

		return nil
	})

	if err != nil {
		return api.ReassignmentSummary{}, err
	}

	return summary, nil
}
