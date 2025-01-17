package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
