package model

import "github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"

type User struct {
	BaseModel
	UserId   string `gorm:"uniqueIndex"`
	IsActive bool
	Username string
	TeamName string
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
