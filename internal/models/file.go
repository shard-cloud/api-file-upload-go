package models

import (
	"time"
	"gorm.io/gorm"
)

type File struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	OriginalName string   `json:"original_name" gorm:"not null"`
	Path        string    `json:"path" gorm:"not null"`
	Size        int64     `json:"size" gorm:"not null"`
	MimeType    string    `json:"mime_type" gorm:"not null"`
	Extension   string    `json:"extension" gorm:"not null"`
	Hash        string    `json:"hash" gorm:"uniqueIndex;not null"`
	UploadedAt  time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (File) TableName() string {
	return "files"
}
