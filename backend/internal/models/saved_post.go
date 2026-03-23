package models

import "time"

type SavedPost struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_saved_post" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_saved_post" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
