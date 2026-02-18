package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/gorm"
)

type CommentHandler struct {
	DB *gorm.DB
}

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

type CommentResponse struct {
	ID                uint              `json:"id"`
	Content           string            `json:"content"`
	VoteCount         int               `json:"vote_count"`
	CreatedAt         time.Time         `json:"created_at"`
	UserID            uint              `json:"user_id"`
	AnonymousUsername string            `json:"anonymous_username"`
	AvatarHash        string            `json:"avatar_hash"`
	PostID            uint              `json:"post_id"`
	ParentID          *uint             `json:"parent_id"`
	UserVote          int               `json:"user_vote"`
	Replies           []CommentResponse `json:"replies,omitempty"`
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Param("id") // Changed from "post_id" to "id"

	var comments []models.Comment
	if err := h.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Replies.User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	var response []CommentResponse
	for _, comment := range comments {
		response = append(response, h.toCommentResponse(comment, c))
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate post exists
	var post models.Post
	if err := h.DB.First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// If parent comment specified, validate it exists
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := h.DB.First(&parentComment, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
			return
		}
	}

	comment := models.Comment{
		Content:   req.Content,
		UserID:    userID.(uint),
		PostID:    req.PostID,
		ParentID:  req.ParentID,
		VoteCount: 0,
	}

	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load user relationship
	h.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, h.toCommentResponse(comment, c))
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentID := c.Param("id")

	var comment models.Comment
	if err := h.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	if err := h.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (h *CommentHandler) toCommentResponse(comment models.Comment, c *gin.Context) CommentResponse {
	userID, _ := c.Get("userID")

	// Get user vote
	var vote models.Vote
	userVote := 0
	if err := h.DB.Where("user_id = ? AND comment_id = ?", userID, comment.ID).First(&vote).Error; err == nil {
		userVote = vote.Value
	}

	response := CommentResponse{
		ID:                comment.ID,
		Content:           comment.Content,
		VoteCount:         comment.VoteCount,
		CreatedAt:         comment.CreatedAt,
		UserID:            comment.UserID,
		AnonymousUsername: comment.User.AnonymousUsername,
		AvatarHash:        comment.User.AvatarHash,
		PostID:            comment.PostID,
		ParentID:          comment.ParentID,
		UserVote:          userVote,
	}

	// Add replies if they exist
	if len(comment.Replies) > 0 {
		for _, reply := range comment.Replies {
			response.Replies = append(response.Replies, h.toCommentResponse(reply, c))
		}
	}

	return response
}
