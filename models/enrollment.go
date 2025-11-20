package models

import "time"

type Enrollment struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID  int       `gorm:"not null;index:idx_student_program,unique" json:"student_id"`
	ProgramID  int       `gorm:"not null;index:idx_student_program,unique" json:"program_id"`
	Status     string    `gorm:"type:varchar(20);default:'enrolled'" json:"status"`
	EnrolledAt time.Time `gorm:"autoCreateTime" json:"enrolled_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Enrollment) TableName() string {
	return "enrollment"
}

type CreateEnrollmentRequest struct {
	StudentID int `json:"student_id"`
	ProgramID int `json:"program_id"`
}

type UpdateEnrollmentRequest struct {
	Status string `json:"status"`
}
