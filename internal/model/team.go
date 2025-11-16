package model

type Team struct {
	BaseModel
	TeamName string `gorm:"uniqueIndex"`
}
