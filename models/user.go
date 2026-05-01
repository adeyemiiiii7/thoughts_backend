package models

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"unique;not null" json:"username"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	Thoughts  []Thought  `json:"thoughts,omitempty"`
	Comments  []Comment  `json:"comments,omitempty"`
	Reactions []Reaction `json:"reactions,omitempty"`
	Followers []Follow   `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
	Following []Follow   `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Interests []Interest `gorm:"many2many:user_interests;" json:"interests"`
}
