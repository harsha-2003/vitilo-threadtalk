package models

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"uniqueIndex;not null" json:"name" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	IconURL     string `json:"icon_url"`
	CreatedBy uint `gorm:"not null;index" json:"created_by"`

	// Relationships
	Posts       []Post            `gorm:"foreignKey:CommunityID" json:"posts,omitempty"`
	Members     []CommunityMember `gorm:"foreignKey:CommunityID" json:"members,omitempty"`
	MemberCount int               `gorm:"-" json:"member_count"` // Computed field
}

type CommunityMember struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;index:idx_user_community" json:"user_id"`
	CommunityID uint      `gorm:"not null;index:idx_user_community" json:"community_id"`
	JoinedAt    time.Time `json:"joined_at"`
	Role string `gorm:"type:varchar(20);default:'member'" json:"role"`

	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Community Community `gorm:"foreignKey:CommunityID" json:"community,omitempty"`
}
