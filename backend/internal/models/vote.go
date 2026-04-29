package models

import (
	"time"
)

type Vote struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uint  `gorm:"not null;index:idx_user_vote" json:"user_id"`
	PostID    *uint `gorm:"index:idx_user_vote" json:"post_id"`
	CommentID *uint `gorm:"index:idx_user_vote" json:"comment_id"`

	Value int `gorm:"not null" json:"value"` // 1 for upvote, -1 for downvote

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
