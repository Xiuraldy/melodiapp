package models

import "time"

type Service struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StartTime string    `json:"start_time" gorm:"column:start_time"`
	EndTime   string    `json:"end_time" gorm:"column:end_time"`
	Name      string    `json:"name"`
	CreatedBy uint      `json:"created_by" gorm:"column:created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
