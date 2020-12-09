package models

import "time"

type BaseUser struct {
	ID        int64     `gorm:"primary_key"`
	UserName  string    `gorm:"not null; type:varchar(50)"`
	Email     string    `gorm:"not null; unique; type:varchar(200)" binding:"required"`
	Password  string    `gorm:"not null; type:varchar(100)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}
