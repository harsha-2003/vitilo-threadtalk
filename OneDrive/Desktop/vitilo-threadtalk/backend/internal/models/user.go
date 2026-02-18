package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email             string `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash      string `gorm:"not null" json:"-"`
	AnonymousUsername string `gorm:"uniqueIndex;not null" json:"anonymous_username"`
	AvatarHash        string `gorm:"not null" json:"avatar_hash"` // For Jdenticon
	IsVerified        bool   `gorm:"default:false" json:"is_verified"`

	// Relationships
	Posts       []Post            `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments    []Comment         `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Votes       []Vote            `gorm:"foreignKey:UserID" json:"votes,omitempty"`
	Communities []CommunityMember `gorm:"foreignKey:UserID" json:"communities,omitempty"`
}
