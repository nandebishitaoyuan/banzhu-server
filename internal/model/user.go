package model

import (
	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `json:"id" gorm:"primaryKey;type:bigint"`
	Username  string         `json:"username" gorm:"not null;uniqueIndex;size:64"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt *JSONTime      `json:"createdAt"`
	UpdatedAt *JSONTime      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
