// models/service_outfit.go (o donde tengas tus modelos)
package models

type ServiceOutfit struct {
	ServiceID uint `json:"service_id" gorm:"primaryKey"`
	OutfitID  uint `json:"outfit_id" gorm:"primaryKey"`
}
