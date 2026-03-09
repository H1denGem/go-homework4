package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64  `gorm:"primaryKey"`
	Content   string `gorm:"not null"`
	UserID    int64  `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	PostID    int64  `gorm:"not null"`
	Post      Post   `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
