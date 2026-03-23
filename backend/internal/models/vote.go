package models

import (
	"time"
)

type Vote struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

UserID    uint  `gorm:"not null;uniqueIndex:idx_user_post_vote;uniqueIndex:idx_user_comment_vote" json:"user_id"`
PostID    *uint `gorm:"uniqueIndex:idx_user_post_vote" json:"post_id"`
CommentID *uint `gorm:"uniqueIndex:idx_user_comment_vote" json:"comment_id"`
	//changed

	Value int `gorm:"not null" json:"value"` // 1 for upvote, -1 for downvote

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
