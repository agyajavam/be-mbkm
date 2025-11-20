package models

import "time"

type Lecturer struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null;uniqueIndex" json:"user_id"`
	NIDN       string    `gorm:"column:nidn;type:varchar(20);uniqueIndex;not null" json:"nidn"`
	FullName   string    `gorm:"type:varchar(100);not null" json:"full_name"`
	Phone      string    `gorm:"type:varchar(20)" json:"phone"`
	Department string    `gorm:"type:varchar(100)" json:"department"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Lecturer) TableName() string {
	return "lecturer"
}

type CreateLecturerRequest struct {
	UserID     int    `json:"user_id" binding:"required"`
	NIDN       string `json:"nidn" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
}

type UpdateLecturerRequest struct {
	NIDN       string `json:"nidn"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	IsActive   *bool  `json:"is_active"`
}
