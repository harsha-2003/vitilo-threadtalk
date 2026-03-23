package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}
func (h *PostHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif allowed"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Max 5MB allowed"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	savePath := "uploads/" + filename

	// Save file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image_url": "/" + savePath,
		"message":   "Image uploaded successfully",
	})
}

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	CommunityID uint   `json:"community_id" binding:"required"`
	PostType    string `json:"post_type"`
	ImageURL    string `json:"image_url"`
}
type UpdatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	PostType string `json:"post_type"`
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

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate community exists
	var community models.Community
	if err := h.DB.First(&community, req.CommunityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	// Set default post type
	if req.PostType == "" {
		req.PostType = "text"
	}
switch req.PostType {
case "text":
	if req.ImageURL != "" || req.LinkURL != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text posts cannot include image_url or link_url"})
		return
	}
case "image":
	if strings.TrimSpace(req.ImageURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image posts require image_url"})
		return
	}
	req.Content = ""
	req.LinkURL = ""
case "link":
	if strings.TrimSpace(req.LinkURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link posts require link_url"})
		return
	}
	req.Content = ""
	req.ImageURL = ""
default:
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_type. Allowed values: text, image, link"})
	return
}
	post := models.Post{
		Title:       req.Title,
		Content:     req.Content,
		PostType:    req.PostType,
		ImageURL:    req.ImageURL,
		UserID:      userID.(uint),
		CommunityID: req.CommunityID,
		VoteCount:   0,
		LinkURL: req.LinkURL,
	}

	if err := h.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Load relationships
	h.DB.Preload("User").Preload("Community").First(&post, post.ID)

	c.JSON(http.StatusCreated, h.toPostResponse(post, userID.(uint)))
}


	// Check ownership
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	if err := h.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	req.LinkURL = strings.TrimSpace(req.LinkURL)

	if req.PostType == "" {
		req.PostType = post.PostType
	}

	switch req.PostType {
	case "text":
		if req.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}
		req.ImageURL = ""
		req.LinkURL = ""
	case "image":
		if req.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image posts require image_url"})
			return
		}
		req.Content = ""
		req.LinkURL = ""
	case "link":
		if req.LinkURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Link posts require link_url"})
			return
		}
		req.Content = ""
		req.ImageURL = ""
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_type"})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.ImageURL = req.ImageURL
	post.LinkURL = req.LinkURL
	post.PostType = req.PostType

	if err := h.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	h.DB.Preload("User").Preload("Community").First(&post, post.ID)
	c.JSON(http.StatusOK, h.toPostResponse(post, userID))
}
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	postID := c.Param("id")

	var post models.Post
	if err := h.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check ownership
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
	if post.ImageURL != "" {
		imagePath := strings.TrimPrefix(post.ImageURL, "/")
		_ = os.Remove(imagePath)
	}

	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Vote{}).Error; err != nil {
		return err
	}
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
		return err
	}
	if err := tx.Delete(&post).Error; err != nil {
		return err
	}
	return nil
})
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
	return
}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
func (h *PostHandler) GetCommunityPosts(c *gin.Context) {
	userID, _ := getUserID(c)
	communityID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort", "new")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.Post{}).
		Where("community_id = ?", communityID).
		Preload("User").
		Preload("Community")

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch community posts"})
		return
	}

	var total int64
	h.DB.Model(&models.Post{}).Where("community_id = ?", communityID).Count(&total)

	resp := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, h.toPostResponse(p, userID))
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":       resp,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}
func (h *PostHandler) toPostResponse(post models.Post, currentUserID uint) PostResponse {
	// Get user vote
	var vote models.Vote
	userVote := 0
	if err := h.DB.Where("user_id = ? AND post_id = ?", currentUserID, post.ID).First(&vote).Error; err == nil {
		userVote = vote.Value
	}

	// Get comment count
	var commentCount int64
	h.DB.Model(&models.Comment{}).Where("post_id = ?", post.ID).Count(&commentCount)

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
		CommunityName:     post.Community.Name,
		UserVote:          userVote,
		LinkURL:           post.LinkURL,
		
	}
}
