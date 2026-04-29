package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	CommunityID uint   `json:"community_id" binding:"required"`
	PostType    string `json:"post_type"`
	ImageURL    string `json:"image_url"`
}

type PostResponse struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	ImageURL          string    `json:"image_url"`
	PostType          string    `json:"post_type"`
	VoteCount         int       `json:"vote_count"`
	CommentCount      int       `json:"comment_count"`
	CreatedAt         time.Time `json:"created_at"`
	UserID            uint      `json:"user_id"`
	AnonymousUsername string    `json:"anonymous_username"`
	AvatarHash        string    `json:"avatar_hash"`
	CommunityID       uint      `json:"community_id"`
	CommunityName     string    `json:"community_name"`
	UserVote          int       `json:"user_vote"`
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

// -----------------------------
// Create Post
// -----------------------------
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Default post type
	if req.PostType == "" {
		req.PostType = "text"
	}

	// Validate community exists
	var community models.Community
	if err := h.DB.First(&community, req.CommunityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	post := models.Post{
		Title:       req.Title,
		Content:     req.Content,
		PostType:    req.PostType,
		ImageURL:    req.ImageURL,
		UserID:      userID,
		CommunityID: req.CommunityID,
		VoteCount:   0,
	}

	if err := h.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Load relationships for response
	if err := h.DB.Preload("User").Preload("Community").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load created post"})
		return
	}

	c.JSON(http.StatusCreated, h.toPostResponse(post, userID))
}

// -----------------------------
// Upload Image
// -----------------------------
func (h *PostHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif allowed"})
		return
	}

	// max 5MB
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Max 5MB allowed"})
		return
	}

	filename := uuid.New().String() + ext
	savePath := "uploads/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image_url": "/" + savePath,
		"message":   "Image uploaded successfully",
	})
}

// -----------------------------
// Get Feed Posts
// -----------------------------


func (h *PostHandler) GetUserPosts(c *gin.Context) {
	currentUserIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUserID := currentUserIDAny.(uint)

	userIDParam := c.Param("id")
	userID64, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	targetUserID := uint(userID64)

	var posts []models.Post
	if err := h.DB.Model(&models.Post{}).
		Where("user_id = ?", targetUserID).
		Preload("User").
		Preload("Community").
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		response = append(response, h.toPostResponse(post, currentUserID))
	}

	c.JSON(http.StatusOK, gin.H{"posts": response})
}
func (h *PostHandler) GetPosts(c *gin.Context) {
	userID, _ := getUserID(c) // userID=0 if missing, safe for user_vote queries

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort", "new")
	communityID := c.Query("community_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := h.DB.Model(&models.Post{}).
		Preload("User").
		Preload("Community")

	// Filter by community
	if communityID != "" {
		query = query.Where("community_id = ?", communityID)
	}

	// Sort
	switch sortBy {
	case "hot":
		query = query.Order("(vote_count * 1.0 / (julianday('now') - julianday(created_at) + 2)) DESC")
	case "top":
		query = query.Order("vote_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	var posts []models.Post
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	resp := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, h.toPostResponse(p, userID))
	}

	var total int64
	countQuery := h.DB.Model(&models.Post{})
	if communityID != "" {
		countQuery = countQuery.Where("community_id = ?", communityID)
	}
	countQuery.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"posts":       resp,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// -----------------------------
// Get Single Post
// -----------------------------
func (h *PostHandler) GetPost(c *gin.Context) {
	userID, _ := getUserID(c)
	postID := c.Param("id")

	var post models.Post
	if err := h.DB.Preload("User").
		Preload("Community").
		First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, h.toPostResponse(post, userID))
}

// -----------------------------
// Delete Post
// -----------------------------
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")

	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	if err := h.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// -----------------------------
// Get My Posts (Profile page)
// GET /api/users/me/posts
// -----------------------------
func (h *PostHandler) GetMyPosts(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var posts []models.Post
	if err := h.DB.Model(&models.Post{}).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Community").
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	resp := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, h.toPostResponse(p, userID))
	}

	c.JSON(http.StatusOK, gin.H{"posts": resp})
}

// -----------------------------
// Helpers
// -----------------------------
func (h *PostHandler) toPostResponse(post models.Post, currentUserID uint) PostResponse {
	// user vote (safe even if currentUserID == 0)
	userVote := 0
	if currentUserID != 0 {
		var vote models.Vote
		if err := h.DB.Where("user_id = ? AND post_id = ?", currentUserID, post.ID).First(&vote).Error; err == nil {
			userVote = vote.Value
		}
	}

	// comment count
	var commentCount int64
	h.DB.Model(&models.Comment{}).Where("post_id = ?", post.ID).Count(&commentCount)

	communityName := ""
	if post.Community.ID != 0 {
		communityName = post.Community.Name
	}

	return PostResponse{
		ID:                post.ID,
		Title:             post.Title,
		Content:           post.Content,
		ImageURL:          post.ImageURL,
		PostType:          post.PostType,
		VoteCount:         post.VoteCount,
		CommentCount:      int(commentCount),
		CreatedAt:         post.CreatedAt,
		UserID:            post.UserID,
		AnonymousUsername: post.User.AnonymousUsername,
		AvatarHash:        post.User.AvatarHash,
		CommunityID:       post.CommunityID,
		CommunityName:     communityName,
		UserVote:          userVote,
	}
}

func getUserID(c *gin.Context) (uint, bool) {
	v, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
