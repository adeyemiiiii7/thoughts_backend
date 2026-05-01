package models

import "time"

type Thought struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	User      User       `json:"user,omitempty"`
	Title     string     `gorm:"not null" json:"title"`
	Content   string     `gorm:"not null" json:"content"`
	Comments  []Comment  `json:"comments,omitempty"`
	Reactions []Reaction `json:"reactions,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
