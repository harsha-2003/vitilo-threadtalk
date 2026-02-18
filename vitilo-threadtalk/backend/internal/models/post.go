package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title    string `gorm:"not null" json:"title" binding:"required"`
	Content  string `gorm:"type:text" json:"content"`
	ImageURL string `json:"image_url"`
	PostType string `gorm:"default:'text'" json:"post_type"` // text, image, link

	UserID      uint `gorm:"not null;index" json:"user_id"`
	CommunityID uint `gorm:"not null;index" json:"community_id"`

	VoteCount int `gorm:"default:0" json:"vote_count"`

	// Relationships
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Community Community `gorm:"foreignKey:CommunityID" json:"community"`
	Comments  []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Votes     []Vote    `gorm:"foreignKey:PostID" json:"votes,omitempty"`
}
