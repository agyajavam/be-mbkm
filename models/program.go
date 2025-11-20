package models

import "time"

type Program struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Credits     int       `gorm:"default:3" json:"credits"`
	Semester    int       `gorm:"not null" json:"semester"`
	LecturerID  int       `gorm:"not null" json:"lecturer_id"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Program) TableName() string {
	return "program"
}

type CreateProgramRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
	Semester    int    `json:"semester"`
	LecturerID  int    `json:"lecturer_id"`
}

type UpdateProgramRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
	Semester    int    `json:"semester"`
	IsActive    *bool  `json:"is_active"`
}
