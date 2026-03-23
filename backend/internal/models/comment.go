package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Content string `gorm:"type:text;not null" json:"content" binding:"required"`

	UserID   uint  `gorm:"not null;index" json:"user_id"`
	PostID   uint  `gorm:"not null;index" json:"post_id"`
	ParentID *uint `gorm:"index" json:"parent_id"` // For nested comments

	VoteCount int `gorm:"default:0" json:"vote_count"`

	// Relationships
	User    User      `gorm:"foreignKey:UserID" json:"user"`
	Post    Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Votes   []Vote    `gorm:"foreignKey:CommentID" json:"votes,omitempty"`
}
