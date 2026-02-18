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
