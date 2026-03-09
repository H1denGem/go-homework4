package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        int64  `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"not null"`
	UserID    int64  `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
