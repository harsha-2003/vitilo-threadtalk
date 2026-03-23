package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/gorm"
)

type VoteHandler struct {
	DB *gorm.DB
}

type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=-1 1"` // -1 for downvote, 1 for upvote
}

func NewVoteHandler(db *gorm.DB) *VoteHandler {
	return &VoteHandler{DB: db}
}

func (h *VoteHandler) VotePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	postID := c.Param("id")

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if post exists
	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check for existing vote
	var existingVote models.Vote
	err := h.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		vote := models.Vote{
			UserID: userID.(uint),
			PostID: &post.ID,
			Value:  req.Value,
		}
		if err := h.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
			return
		}

		// Update post vote count
		h.DB.Model(&post).UpdateColumn("vote_count", gorm.Expr("vote_count + ?", req.Value))

	} else if err == nil {
		// Update existing vote
		if existingVote.Value == req.Value {
			// Remove vote if same value
			h.DB.Delete(&existingVote)
			h.DB.Model(&post).UpdateColumn("vote_count", gorm.Expr("vote_count - ?", req.Value))
		} else {
			// Change vote
			diff := req.Value - existingVote.Value
			h.DB.Model(&existingVote).Update("value", req.Value)
			h.DB.Model(&post).UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff))
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
		return
	}

	// Get updated post
	h.DB.First(&post, postID)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Vote processed successfully",
		"vote_count": post.VoteCount,
	})
}

func (h *VoteHandler) VoteComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentID := c.Param("id")

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if comment exists
	var comment models.Comment
	if err := h.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check for existing vote
	var existingVote models.Vote
	err := h.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingVote).Error

	if err == gorm.ErrRecordNotFound {
		// Create new vote
		vote := models.Vote{
			UserID:    userID.(uint),
			CommentID: &comment.ID,
			Value:     req.Value,
		}
		if err := h.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
			return
		}

		// Update comment vote count
		h.DB.Model(&comment).UpdateColumn("vote_count", gorm.Expr("vote_count + ?", req.Value))

	} else if err == nil {
		// Update existing vote
		if existingVote.Value == req.Value {
			// Remove vote if same value
			h.DB.Delete(&existingVote)
			h.DB.Model(&comment).UpdateColumn("vote_count", gorm.Expr("vote_count - ?", req.Value))
		} else {
			// Change vote
			diff := req.Value - existingVote.Value
			h.DB.Model(&existingVote).Update("value", req.Value)
			h.DB.Model(&comment).UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff))
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
		return
	}

	// Get updated comment
	h.DB.First(&comment, commentID)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Vote processed successfully",
		"vote_count": comment.VoteCount,
	})
}
