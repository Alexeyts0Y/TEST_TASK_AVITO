package model

import (
	"strings"
	"time"

	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId   string `gorm:"uniqueIndex"`
	IsActive bool
	Username string
	TeamName string
}

type Team struct {
	gorm.Model
	TeamName string `gorm:"uniqueIndex"`
}

type PullRequest struct {
	gorm.Model
	AssignedReviewers string
	AuthorId          string
	CreatedAt         *time.Time
	MergedAt          *time.Time
	PullRequestId     string `gorm:"uniqueIndex"`
	PullRequestName   string
	Status            api.PullRequestStatus
}

func (u *User) ToAPIUser() api.User {
	return api.User{
		UserId:   u.UserId,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func (u *User) ToAPITeamMember() api.TeamMember {
	return api.TeamMember{
		UserId:   u.UserId,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}

func (pr *PullRequest) ToAPIPullRequest() api.PullRequest {
	reviewers := []string{}
	if pr.AssignedReviewers != "" {
		reviewers = strings.Split(pr.AssignedReviewers, ",")
	}

	return api.PullRequest{
		AssignedReviewers: reviewers,
		AuthorId:          pr.AuthorId,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            pr.Status,
	}
}

func FromAPIPullRequest(apiPr api.PullRequest) PullRequest {
	reviewers := strings.Join(apiPr.AssignedReviewers, ",")

	return PullRequest{
		AssignedReviewers: reviewers,
		AuthorId:          apiPr.AuthorId,
		CreatedAt:         apiPr.CreatedAt,
		MergedAt:          apiPr.MergedAt,
		PullRequestId:     apiPr.PullRequestId,
		PullRequestName:   apiPr.PullRequestName,
		Status:            apiPr.Status,
	}
}
