package models

type ServiceSong struct {
	ServiceID uint `json:"service_id" gorm:"primaryKey;column:service_id"`
	SongID    uint `json:"song_id" gorm:"primaryKey;column:song_id"`
}
