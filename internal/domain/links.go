package domain

import "time"

type Link struct {
	ID             int64     `json:"id" db:"id"`
	Description    string    `json:"description" db:"description"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	IsOnceDownload string    `json:"is_once_download" db:"is_once_download"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

func (Link) TableName() string { return "links" }

type LinksAnalytics struct {
	LinkID         int64     `json:"link_id" db:"link_id"`
	Downloads      int       `json:"downloads" db:"downloads"`
	LastDownloadAt time.Time `json:"last_download_at" db:"last_download_at"`
}

func (LinksAnalytics) TableName() string { return "links_analytics" }
