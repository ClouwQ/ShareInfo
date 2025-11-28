package domain

import "time"

type UploadedFile struct {
	LinkId    int64     `json:"link_id" db:"link_id"`
	Name      string    `json:"name" db:"name"`
	Size      int64     `json:"size" db:"size"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (UploadedFile) TableName() string { return "uploaded_files" }
