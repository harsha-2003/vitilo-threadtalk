package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/gorm"
)

type CommunityHandler struct {
	DB *gorm.DB
}

type CreateCommunityRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
}

type CommunityResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	MemberCount int       `json:"member_count"`
	IsMember    bool      `json:"is_member"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewCommunityHandler(db *gorm.DB) *CommunityHandler {
	return &CommunityHandler{DB: db}
}

func (h *CommunityHandler) CreateCommunity(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if community name already exists
	var existingCommunity models.Community
	if err := h.DB.Where("name = ?", req.Name).First(&existingCommunity).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Community name already exists"})
		return
	}

	community := models.Community{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
	}

	if err := h.DB.Create(&community).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create community"})
		return
	}

	// Auto-join the creator
	member := models.CommunityMember{
		UserID:      userID.(uint),
		CommunityID: community.ID,
		JoinedAt:    time.Now(),
	}
	h.DB.Create(&member)

	c.JSON(http.StatusCreated, h.toCommunityResponse(community, userID.(uint)))
}

func (h *CommunityHandler) GetCommunities(c *gin.Context) {
	userID, _ := c.Get("userID")

	var communities []models.Community
	if err := h.DB.Find(&communities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch communities"})
		return
	}

	var response []CommunityResponse
	for _, community := range communities {
		response = append(response, h.toCommunityResponse(community, userID.(uint)))
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommunityHandler) GetCommunity(c *gin.Context) {
	userID, _ := c.Get("userID")
	communityID := c.Param("id")

	var community models.Community
	if err := h.DB.First(&community, communityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	c.JSON(http.StatusOK, h.toCommunityResponse(community, userID.(uint)))
}

func (h *CommunityHandler) JoinCommunity(c *gin.Context) {
	userID, _ := c.Get("userID")
	communityID := c.Param("id")

	// Check if community exists
	var community models.Community
	if err := h.DB.First(&community, communityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	// Check if already a member
	var existingMember models.CommunityMember
	if err := h.DB.Where("user_id = ? AND community_id = ?", userID, communityID).
		First(&existingMember).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already a member"})
		return
	}

	member := models.CommunityMember{
		UserID:      userID.(uint),
		CommunityID: community.ID,
		JoinedAt:    time.Now(),
	}

	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join community"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined community"})
}

func (h *CommunityHandler) LeaveCommunity(c *gin.Context) {
	userID, _ := c.Get("userID")
	communityID := c.Param("id")

	result := h.DB.Where("user_id = ? AND community_id = ?", userID, communityID).
		Delete(&models.CommunityMember{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave community"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not a member of this community"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left community"})
}

func (h *CommunityHandler) GetUserCommunities(c *gin.Context) {
	userID, _ := c.Get("userID")

	var members []models.CommunityMember
	if err := h.DB.Where("user_id = ?", userID).
		Preload("Community").
		Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch communities"})
		return
	}

	var response []CommunityResponse
	for _, member := range members {
		response = append(response, h.toCommunityResponse(member.Community, userID.(uint)))
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommunityHandler) toCommunityResponse(community models.Community, currentUserID uint) CommunityResponse {
	// Get member count
	var memberCount int64
	h.DB.Model(&models.CommunityMember{}).Where("community_id = ?", community.ID).Count(&memberCount)

	// Check if current user is a member
	var membership models.CommunityMember
	isMember := false
	if err := h.DB.Where("user_id = ? AND community_id = ?", currentUserID, community.ID).
		First(&membership).Error; err == nil {
		isMember = true
	}

	return CommunityResponse{
		ID:          community.ID,
		Name:        community.Name,
		Description: community.Description,
		IconURL:     community.IconURL,
		MemberCount: int(memberCount),
		IsMember:    isMember,
		CreatedAt:   community.CreatedAt,
	}
}
