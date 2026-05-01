package models

import "time"

type Follow struct {
	FollowerID  uint      `gorm:"primaryKey;autoIncrement:false" json:"follower_id"`
	FollowingID uint      `gorm:"primaryKey;autoIncrement:false" json:"following_id"`
	Follower    User      `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE;" json:"follower,omitempty"`
	Following   User      `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE;" json:"following,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
