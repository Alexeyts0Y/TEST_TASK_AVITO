package repository

import (
	"gorm.io/gorm"
)

type Repository interface {
	TeamRepository
	UserRepository
	PullRequestRepository
	StatsRepository
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}
