package models

import (
	"gorm.io/gorm"
	"time"
)

type GormModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Project struct {
	GormModel
	Name string `json:"name"`
}
