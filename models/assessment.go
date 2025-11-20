package models

import "time"

type Assessment struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	EnrollmentID int       `gorm:"not null" json:"enrollment_id"`
	StudentID    int       `gorm:"not null" json:"student_id"`
	ProgramID    int       `gorm:"not null" json:"program_id"`
	Category     string    `gorm:"type:varchar(50);not null" json:"category"`
	Score        float64   `gorm:"type:decimal(5,2);default:0" json:"score"`
	MaxScore     float64   `gorm:"type:decimal(5,2)" json:"max_score"`
	Weight       float64   `gorm:"type:decimal(5,2);default:0" json:"weight"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Assessment) TableName() string {
	return "assessment"
}

type CreateAssessmentRequest struct {
	EnrollmentID int     `json:"enrollment_id"`
	Category     string  `json:"category"`
	Score        float64 `json:"score"`
	MaxScore     float64 `json:"max_score"`
	Weight       float64 `json:"weight"`
	Notes        string  `json:"notes"`
}

type UpdateAssessmentRequest struct {
	Score    float64 `json:"score"`
	MaxScore float64 `json:"max_score"`
	Weight   float64 `json:"weight"`
	Notes    string  `json:"notes"`
}
