package models

type ServiceUser struct {
	ServiceID uint   `json:"service_id" gorm:"primaryKey;column:service_id"`
	UserID    uint   `json:"user_id" gorm:"primaryKey;column:user_id"`
	Status    string `json:"status"`
}
