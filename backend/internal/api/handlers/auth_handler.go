package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID                uint   `json:"id"`
	Email             string `json:"email"`
	AnonymousUsername string `json:"anonymous_username"`
	AvatarHash        string `json:"avatar_hash"`
}

type ProfileSummaryResponse struct {
	ID                uint   `json:"id"`
	Email             string `json:"email"`
	AnonymousUsername string `json:"anonymous_username"`
	AvatarHash        string `json:"avatar_hash"`
	PostCount         int64  `json:"post_count"`
	CommentCount      int64  `json:"comment_count"`
	CommunityCount    int64  `json:"community_count"`
	Karma             int64  `json:"karma"`
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate UF email
	ufDomain := os.Getenv("UF_EMAIL_DOMAIN")
	if !strings.HasSuffix(strings.ToLower(req.Email), "@"+ufDomain) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must use a UF email address"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate anonymous username and avatar hash
	anonUsername := generateAnonymousUsername()
	avatarHash := generateAvatarHash()

	// Create user
	user := models.User{
		Email:             req.Email,
		PasswordHash:      string(hashedPassword),
		AnonymousUsername: anonUsername,
		AvatarHash:        avatarHash,
		IsVerified:        false,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:                user.ID,
			Email:             user.Email,
			AnonymousUsername: user.AnonymousUsername,
			AvatarHash:        user.AvatarHash,
		},
	})
}
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var postCount, commentCount, communityCount int64
	var postVoteSum, commentVoteSum int64

	h.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
	h.DB.Model(&models.Comment{}).Where("user_id = ?", user.ID).Count(&commentCount)
	h.DB.Model(&models.CommunityMember{}).Where("user_id = ?", user.ID).Count(&communityCount)
	h.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Select("COALESCE(SUM(vote_count), 0)").Scan(&postVoteSum)
	h.DB.Model(&models.Comment{}).Where("user_id = ?", user.ID).Select("COALESCE(SUM(vote_count), 0)").Scan(&commentVoteSum)

	c.JSON(http.StatusOK, gin.H{
		"id":                 user.ID,
		"anonymous_username": user.AnonymousUsername,
		"avatar_hash":        user.AvatarHash,
		"created_at":         user.CreatedAt,
		"post_count":         postCount,
		"comment_count":      commentCount,
		"community_count":    communityCount,
		"karma":              postVoteSum + commentVoteSum,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:                user.ID,
			Email:             user.Email,
			AnonymousUsername: user.AnonymousUsername,
			AvatarHash:        user.AvatarHash,
		},
	})
}
func (h *AuthHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var postCount, commentCount, communityCount int64
	var postVoteSum, commentVoteSum int64

	h.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&postCount)
	h.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	h.DB.Model(&models.CommunityMember{}).Where("user_id = ?", userID).Count(&communityCount)
	h.DB.Model(&models.Post{}).Where("user_id = ?", userID).Select("COALESCE(SUM(vote_count), 0)").Scan(&postVoteSum)
	h.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Select("COALESCE(SUM(vote_count), 0)").Scan(&commentVoteSum)

	c.JSON(http.StatusOK, ProfileSummaryResponse{
		ID:                user.ID,
		Email:             user.Email,
		AnonymousUsername: user.AnonymousUsername,
		AvatarHash:        user.AvatarHash,
		PostCount:         postCount,
		CommentCount:      commentCount,
		CommunityCount:    communityCount,
		Karma:             postVoteSum + commentVoteSum,
	})
}
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:                user.ID,
		Email:             user.Email,
		AnonymousUsername: user.AnonymousUsername,
		AvatarHash:        user.AvatarHash,
	})
}

// Helper functions
func generateAnonymousUsername() string {
	adjectives := []string{"Happy", "Silent", "Brave", "Quick", "Gentle", "Proud", "Clever", "Calm"}
	nouns := []string{"Gator", "Tiger", "Eagle", "Falcon", "Panther", "Wolf", "Bear", "Fox"}

	randomBytes := make([]byte, 2)
	rand.Read(randomBytes)

	adjIdx := int(randomBytes[0]) % len(adjectives)
	nounIdx := int(randomBytes[1]) % len(nouns)

	return fmt.Sprintf("%s%s%d", adjectives[adjIdx], nouns[nounIdx], randomBytes[0]%100)
}

func generateAvatarHash() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}

func generateJWTToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
