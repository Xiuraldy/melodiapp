package models

import "time"

type Song struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	SongKey       string    `json:"song_key" gorm:"column:song_key"`
	BPM           int       `json:"bpm"`
	TimeSignature string    `json:"time_signature" gorm:"column:time_signature"`
	Duration      string    `json:"duration"`
	Structure     string    `json:"structure" gorm:"column:structure"`
	HasSequence   string    `json:"has_sequence" gorm:"column:has_sequence"`
	HasChart      string    `json:"has_chart" gorm:"column:has_chart"`
	HasScore      string    `json:"has_score" gorm:"column:has_score"`
	YoutubeURL    string    `json:"youtube_url" gorm:"column:youtube_url"`
	VoiceURL      string    `json:"voice_url" gorm:"column:voice_url"`
	GuitarURL     string    `json:"guitar_url" gorm:"column:guitar_url"`
	PianoURL      string    `json:"piano_url" gorm:"column:piano_url"`
	DrumsURL      string    `json:"drums_url" gorm:"column:drums_url"`
	BassURL       string    `json:"bass_url" gorm:"column:bass_url"`
	WindURL       string    `json:"wind_url" gorm:"column:wind_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
