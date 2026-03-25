package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ==================== Test Setup ====================

var testUserCounter = 0

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // ✅ DISABLE LOGS
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Community{}, &models.CommunityMember{}, &models.Post{}, &models.Comment{}, &models.Vote{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("UF_EMAIL_DOMAIN", "ufl.edu")
	return gin.Default()
}

func createUser(t *testing.T, db *gorm.DB, email string) *models.User {
	testUserCounter++
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := models.User{
		Email:             email,
		PasswordHash:      string(hashedPassword),
		AnonymousUsername: fmt.Sprintf("User_%d_%s", testUserCounter, email[:5]),
		AvatarHash:        fmt.Sprintf("hash_%d", testUserCounter),
		IsVerified:        false,
	}
	db.Create(&user)
	return &user
}

func createCommunity(t *testing.T, db *gorm.DB, name string) *models.Community {
	comm := models.Community{
		Name:        name,
		Description: "Test community",
		IconURL:     "https://example.com/icon.png",
	}
	db.Create(&comm)
	return &comm
}

func createPost(t *testing.T, db *gorm.DB, title string, userID, communityID uint) *models.Post {
	post := models.Post{
		Title:       title,
		Content:     "Test content",
		PostType:    "text",
		UserID:      userID,
		CommunityID: communityID,
		VoteCount:   0,
	}
	db.Create(&post)
	return &post
}

func createComment(t *testing.T, db *gorm.DB, content string, userID, postID uint) *models.Comment {
	comment := models.Comment{
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		VoteCount: 0,
	}
	db.Create(&comment)
	return &comment
}

func performRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ==================== Auth Handler Tests ====================

// TEST 1: Register Success
func TestAuthRegisterSuccess(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "student1@ufl.edu",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/register", req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

// TEST 2: Register - Invalid Email
func TestAuthRegisterInvalidEmail(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "not-an-email",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/register", req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// TEST 3: Register - Non UF Email
func TestAuthRegisterNonUFEmail(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "user@gmail.com",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/register", req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// TEST 4: Register - Short Password
func TestAuthRegisterShortPassword(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "student2@ufl.edu",
		Password: "short",
	}

	w := performRequest(router, "POST", "/register", req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// TEST 5: Register - Duplicate Email
func TestAuthRegisterDuplicateEmail(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)

	req := RegisterRequest{
		Email:    "duplicate@ufl.edu",
		Password: "password123",
	}

	performRequest(router, "POST", "/register", req)
	w := performRequest(router, "POST", "/register", req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, got %d", w.Code)
	}
}

// TEST 6: Login Success
func TestAuthLoginSuccess(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	registerReq := RegisterRequest{
		Email:    "login@ufl.edu",
		Password: "password123",
	}
	performRequest(router, "POST", "/register", registerReq)

	loginReq := LoginRequest{
		Email:    "login@ufl.edu",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/login", loginReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 7: Login - Wrong Password
func TestAuthLoginWrongPassword(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	registerReq := RegisterRequest{
		Email:    "wrongpass@ufl.edu",
		Password: "password123",
	}
	performRequest(router, "POST", "/register", registerReq)

	loginReq := LoginRequest{
		Email:    "wrongpass@ufl.edu",
		Password: "wrongpassword",
	}

	w := performRequest(router, "POST", "/login", loginReq)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// TEST 8: Login - User Not Found
func TestAuthLoginUserNotFound(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	router.POST("/login", handler.Login)

	loginReq := LoginRequest{
		Email:    "notexist@ufl.edu",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/login", loginReq)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// TEST 9: Get Current User
func TestAuthGetCurrentUser(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)
	user := createUser(t, db, "current@ufl.edu")

	router.GET("/me", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetCurrentUser(c)
	})

	w := performRequest(router, "GET", "/me", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 10: Get Current User - Not Found
func TestAuthGetCurrentUserNotFound(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewAuthHandler(db)

	router.GET("/me", func(c *gin.Context) {
		c.Set("userID", uint(9999))
		handler.GetCurrentUser(c)
	})

	w := performRequest(router, "GET", "/me", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TEST 11: Create Community
func TestCreateCommunity(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommunityHandler(db)
	user := createUser(t, db, "creator@ufl.edu")

	router.POST("/communities", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreateCommunity(c)
	})

	req := CreateCommunityRequest{
		Name:        "test-community",
		Description: "A test community",
	}

	w := performRequest(router, "POST", "/communities", req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

// TEST 12: Get Communities
func TestGetCommunities(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommunityHandler(db)
	user := createUser(t, db, "getter@ufl.edu")

	createCommunity(t, db, "comm1")
	createCommunity(t, db, "comm2")

	router.GET("/communities", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetCommunities(c)
	})

	w := performRequest(router, "GET", "/communities", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 13: Join Community
func TestJoinCommunity(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommunityHandler(db)
	user := createUser(t, db, "member@ufl.edu")
	comm := createCommunity(t, db, "joinable")

	router.POST("/communities/:id/join", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.JoinCommunity(c)
	})

	w := performRequest(router, "POST", fmt.Sprintf("/communities/%d/join", comm.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 14: Leave Community
func TestLeaveCommunity(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommunityHandler(db)
	user := createUser(t, db, "leaver@ufl.edu")
	comm := createCommunity(t, db, "leaveable")

	member := models.CommunityMember{
		UserID:      user.ID,
		CommunityID: comm.ID,
	}
	db.Create(&member)

	router.POST("/communities/:id/leave", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.LeaveCommunity(c)
	})

	w := performRequest(router, "POST", fmt.Sprintf("/communities/%d/leave", comm.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 15: Get Community
func TestGetCommunity(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommunityHandler(db)
	user := createUser(t, db, "communityuser@ufl.edu")
	comm := createCommunity(t, db, "getcommunity")

	router.GET("/communities/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetCommunity(c)
	})

	w := performRequest(router, "GET", fmt.Sprintf("/communities/%d", comm.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 16: Create Post
func TestCreatePost(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "poster@ufl.edu")
	comm := createCommunity(t, db, "post-comm")

	router.POST("/posts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreatePost(c)
	})

	req := CreatePostRequest{
		Title:       "Test Post",
		Content:     "Test content",
		CommunityID: comm.ID,
		PostType:    "text",
	}

	w := performRequest(router, "POST", "/posts", req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

// TEST 17: Create Post - Empty Title
func TestCreatePostEmptyTitle(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "poster2@ufl.edu")
	comm := createCommunity(t, db, "post-comm2")

	router.POST("/posts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreatePost(c)
	})

	req := CreatePostRequest{
		Title:       "",
		Content:     "Content",
		CommunityID: comm.ID,
	}

	w := performRequest(router, "POST", "/posts", req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// TEST 18: Create Post - Invalid Community
func TestCreatePostInvalidCommunity(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "poster3@ufl.edu")

	router.POST("/posts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreatePost(c)
	})

	req := CreatePostRequest{
		Title:       "Test",
		Content:     "Content",
		CommunityID: 9999,
	}

	w := performRequest(router, "POST", "/posts", req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TEST 19: Get Posts
func TestGetPosts(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "getter-posts@ufl.edu")
	comm := createCommunity(t, db, "get-posts-comm")

	createPost(t, db, "Post 1", user.ID, comm.ID)
	createPost(t, db, "Post 2", user.ID, comm.ID)

	router.GET("/posts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetPosts(c)
	})

	w := performRequest(router, "GET", "/posts", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 20: Get Single Post
func TestGetPost(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "getter-post@ufl.edu")
	comm := createCommunity(t, db, "get-post-comm")
	post := createPost(t, db, "Single Post", user.ID, comm.ID)

	router.GET("/posts/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetPost(c)
	})

	w := performRequest(router, "GET", fmt.Sprintf("/posts/%d", post.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 21: Delete Post
func TestDeletePost(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "deleter@ufl.edu")
	comm := createCommunity(t, db, "delete-comm")
	post := createPost(t, db, "Delete Me", user.ID, comm.ID)

	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.DeletePost(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/posts/%d", post.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 22: Delete Post - Forbidden
func TestDeletePostForbidden(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user1 := createUser(t, db, "owner@ufl.edu")
	user2 := createUser(t, db, "other@ufl.edu")
	comm := createCommunity(t, db, "forbidden-comm")
	post := createPost(t, db, "Protected", user1.ID, comm.ID)

	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("userID", user2.ID)
		handler.DeletePost(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/posts/%d", post.ID), nil)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", w.Code)
	}
}

// TEST 23: Get User Posts
func TestGetUserPosts(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "userposts@ufl.edu")
	comm := createCommunity(t, db, "user-posts-comm")

	createPost(t, db, "User Post 1", user.ID, comm.ID)
	createPost(t, db, "User Post 2", user.ID, comm.ID)

	router.GET("/users/:id/posts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetUserPosts(c)
	})

	w := performRequest(router, "GET", fmt.Sprintf("/users/%d/posts", user.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 24: Get My Posts
func TestGetMyPosts(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewPostHandler(db)
	user := createUser(t, db, "myposts@ufl.edu")
	comm := createCommunity(t, db, "my-posts-comm")

	createPost(t, db, "My Post 1", user.ID, comm.ID)
	createPost(t, db, "My Post 2", user.ID, comm.ID)

	router.GET("/myposts", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetMyPosts(c)
	})

	w := performRequest(router, "GET", "/myposts", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 25: Create Comment
func TestCreateComment(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommentHandler(db)
	user := createUser(t, db, "commenter@ufl.edu")
	comm := createCommunity(t, db, "comment-comm")
	post := createPost(t, db, "Comment Post", user.ID, comm.ID)

	router.POST("/comments", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreateComment(c)
	})

	req := CreateCommentRequest{
		Content: "Great post!",
		PostID:  post.ID,
	}

	w := performRequest(router, "POST", "/comments", req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

// TEST 26: Create Comment - Invalid Post
func TestCreateCommentInvalidPost(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommentHandler(db)
	user := createUser(t, db, "commenter2@ufl.edu")

	router.POST("/comments", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.CreateComment(c)
	})

	req := CreateCommentRequest{
		Content: "Comment",
		PostID:  9999,
	}

	w := performRequest(router, "POST", "/comments", req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TEST 27: Get Comments
func TestGetComments(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommentHandler(db)
	user := createUser(t, db, "commentgetter@ufl.edu")
	comm := createCommunity(t, db, "comments-comm")
	post := createPost(t, db, "Comments Post", user.ID, comm.ID)

	createComment(t, db, "Comment 1", user.ID, post.ID)
	createComment(t, db, "Comment 2", user.ID, post.ID)

	router.GET("/posts/:id/comments", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.GetComments(c)
	})

	w := performRequest(router, "GET", fmt.Sprintf("/posts/%d/comments", post.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 28: Delete Comment
func TestDeleteComment(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommentHandler(db)
	user := createUser(t, db, "commenter-deleter@ufl.edu")
	comm := createCommunity(t, db, "comment-delete-comm")
	post := createPost(t, db, "Delete Comment", user.ID, comm.ID)
	comment := createComment(t, db, "Delete me", user.ID, post.ID)

	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.DeleteComment(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/comments/%d", comment.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 29: Delete Comment - Forbidden
func TestDeleteCommentForbidden(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewCommentHandler(db)
	user1 := createUser(t, db, "commenter-owner@ufl.edu")
	user2 := createUser(t, db, "commenter-other@ufl.edu")
	comm := createCommunity(t, db, "comment-forbidden-comm")
	post := createPost(t, db, "Forbidden Comment", user1.ID, comm.ID)
	comment := createComment(t, db, "Protected", user1.ID, post.ID)

	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("userID", user2.ID)
		handler.DeleteComment(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/comments/%d", comment.ID), nil)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", w.Code)
	}
}

// TEST 30: Vote Post Upvote
func TestVotePostUpvote(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewVoteHandler(db)
	user := createUser(t, db, "voter@ufl.edu")
	comm := createCommunity(t, db, "vote-comm")
	post := createPost(t, db, "Vote Post", user.ID, comm.ID)

	router.POST("/posts/:id/vote", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.VotePost(c)
	})

	req := VoteRequest{Value: 1}

	w := performRequest(router, "POST", fmt.Sprintf("/posts/%d/vote", post.ID), req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 31: Vote Post Downvote
func TestVotePostDownvote(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewVoteHandler(db)
	user := createUser(t, db, "downvoter@ufl.edu")
	comm := createCommunity(t, db, "downvote-comm")
	post := createPost(t, db, "Downvote Post", user.ID, comm.ID)

	router.POST("/posts/:id/vote", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.VotePost(c)
	})

	req := VoteRequest{Value: -1}

	w := performRequest(router, "POST", fmt.Sprintf("/posts/%d/vote", post.ID), req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 32: Vote Post - Invalid Value
func TestVotePostInvalidValue(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewVoteHandler(db)
	user := createUser(t, db, "invalidvoter@ufl.edu")
	comm := createCommunity(t, db, "invalid-vote-comm")
	post := createPost(t, db, "Invalid Vote Post", user.ID, comm.ID)

	router.POST("/posts/:id/vote", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.VotePost(c)
	})

	req := VoteRequest{Value: 5}

	w := performRequest(router, "POST", fmt.Sprintf("/posts/%d/vote", post.ID), req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// TEST 33: Vote Comment
func TestVoteComment(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewVoteHandler(db)
	user := createUser(t, db, "commentvoter@ufl.edu")
	comm := createCommunity(t, db, "comment-vote-comm")
	post := createPost(t, db, "Comment Vote Post", user.ID, comm.ID)
	comment := createComment(t, db, "Vote comment", user.ID, post.ID)

	router.POST("/comments/:id/vote", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.VoteComment(c)
	})

	req := VoteRequest{Value: 1}

	w := performRequest(router, "POST", fmt.Sprintf("/comments/%d/vote", comment.ID), req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TEST 34: Vote Comment - Invalid Comment
func TestVoteCommentInvalidComment(t *testing.T) {
	db := setupDB(t)
	router := setupRouter()
	handler := NewVoteHandler(db)
	user := createUser(t, db, "invalidcommentvoter@ufl.edu")

	router.POST("/comments/:id/vote", func(c *gin.Context) {
		c.Set("userID", user.ID)
		handler.VoteComment(c)
	})

	req := VoteRequest{Value: 1}

	w := performRequest(router, "POST", "/comments/9999/vote", req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TEST 35: Helper - Generate Avatar Hash
func TestGenerateAvatarHash(t *testing.T) {
	hash := generateAvatarHash()
	if hash == "" {
		t.Error("Hash should not be empty")
	}
	if len(hash) != 32 {
		t.Errorf("Hash length should be 32, got %d", len(hash))
	}
}
