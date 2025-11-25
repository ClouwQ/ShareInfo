package domain

import "time"

type LinkMessages struct {
	LinkID int64     `json:"link_id" db:"link_id"`
	Text   string    `json:"text" db:"text"`
	SentAt time.Time `json:"sent_at" db:"sent_at"`
}

func (LinkMessages) TableName() string { return "link_messages" }
