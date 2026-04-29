//

package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testEnv struct {
	router *gin.Engine
	db     *gorm.DB
}

type authPayload struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("UF_EMAIL_DOMAIN", "ufl.edu")
	_ = os.MkdirAll("uploads", 0755)

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite memory db: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Community{},
		&models.Post{},
		&models.Comment{},
		&models.Vote{},
		&models.CommunityMember{},
	)
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	r := gin.Default()
	SetupRoutes(r, db)
	return &testEnv{router: r, db: db}
}

func performJSONRequest(r http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createMultipartUploadRequest(t *testing.T, path, fieldName, fileName, content string, token string) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create multipart file: %v", err)
	}
	_, _ = part.Write([]byte(content))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func registerAndGetToken(t *testing.T, env *testEnv, email string) (string, uint) {
	t.Helper()
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    email,
		"password": "password123",
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
	}
	var resp authPayload
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}
	return resp.Token, resp.User.ID
}

func createCommunityWithToken(t *testing.T, env *testEnv, token, name string) uint {
	t.Helper()
	w := performJSONRequest(env.router, http.MethodPost, "/api/communities", map[string]any{
		"name":        name,
		"description": "community description",
		"icon_url":    "https://example.com/icon.png",
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create community failed: status=%d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return uint(resp["id"].(float64))
}

func createPostWithToken(t *testing.T, env *testEnv, token string, communityID uint, title string) uint {
	t.Helper()
	w := performJSONRequest(env.router, http.MethodPost, "/api/posts", map[string]any{
		"title":        title,
		"content":      "post content",
		"community_id": communityID,
		"post_type":    "text",
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create post failed: status=%d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return uint(resp["id"].(float64))
}

func createCommentWithToken(t *testing.T, env *testEnv, token string, postID uint, parentID *uint, content string) uint {
	t.Helper()
	payload := map[string]any{
		"content": content,
		"post_id": postID,
	}
	if parentID != nil {
		payload["parent_id"] = *parentID
	}
	w := performJSONRequest(env.router, http.MethodPost, "/api/comments", payload, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment failed: status=%d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return uint(resp["id"].(float64))
}

// ==================== ORIGINAL TESTS ====================

func TestHealthCheck_OK(t *testing.T) {
	env := setupTestEnv(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRegister_Success(t *testing.T) {
	env := setupTestEnv(t)
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@ufl.edu",
		"password": "password123",
	}, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidEmailFormat_ReturnsBadRequest(t *testing.T) {
	env := setupTestEnv(t)
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "not-an-email",
		"password": "password123",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_NonUFEmail_ReturnsBadRequest(t *testing.T) {
	env := setupTestEnv(t)
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@gmail.com",
		"password": "password123",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_DuplicateEmail_ReturnsConflict(t *testing.T) {
	env := setupTestEnv(t)
	_, _ = registerAndGetToken(t, env, "duplicate@ufl.edu")
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "duplicate@ufl.edu",
		"password": "password123",
	}, "")

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_Success(t *testing.T) {
	env := setupTestEnv(t)
	_, _ = registerAndGetToken(t, env, "login@ufl.edu")
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "login@ufl.edu",
		"password": "password123",
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_WrongPassword_ReturnsUnauthorized(t *testing.T) {
	env := setupTestEnv(t)
	_, _ = registerAndGetToken(t, env, "wrongpass@ufl.edu")
	w := performJSONRequest(env.router, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "wrongpass@ufl.edu",
		"password": "wrong-password",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetCurrentUser_WithoutToken_ReturnsUnauthorized(t *testing.T) {
	env := setupTestEnv(t)
	w := performJSONRequest(env.router, http.MethodGet, "/api/auth/me", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetCurrentUser_WithValidToken_ReturnsOK(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "currentuser@ufl.edu")
	w := performJSONRequest(env.router, http.MethodGet, "/api/auth/me", nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateCommunity_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "community1@ufl.edu")
	w := performJSONRequest(env.router, http.MethodPost, "/api/communities", map[string]any{
		"name":        "uf-tech",
		"description": "UF tech club",
	}, token)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateCommunity_DuplicateName_ReturnsConflict(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "community2@ufl.edu")
	_ = createCommunityWithToken(t, env, token, "gators")
	w := performJSONRequest(env.router, http.MethodPost, "/api/communities", map[string]any{
		"name":        "gators",
		"description": "duplicate",
	}, token)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestJoinCommunity_Success(t *testing.T) {
	env := setupTestEnv(t)
	ownerToken, _ := registerAndGetToken(t, env, "owner@ufl.edu")
	memberToken, _ := registerAndGetToken(t, env, "member@ufl.edu")
	communityID := createCommunityWithToken(t, env, ownerToken, "joinable")

	w := performJSONRequest(env.router, http.MethodPost, "/api/communities/"+strconv.Itoa(int(communityID))+"/join", nil, memberToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestJoinCommunity_AlreadyMember_ReturnsConflict(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "member2@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "already-member")

	w := performJSONRequest(env.router, http.MethodPost, "/api/communities/"+strconv.Itoa(int(communityID))+"/join", nil, token)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestLeaveCommunity_NotMember_ReturnsNotFound(t *testing.T) {
	env := setupTestEnv(t)
	ownerToken, _ := registerAndGetToken(t, env, "owner2@ufl.edu")
	strangerToken, _ := registerAndGetToken(t, env, "stranger@ufl.edu")
	communityID := createCommunityWithToken(t, env, ownerToken, "leave-me")

	w := performJSONRequest(env.router, http.MethodPost, "/api/communities/"+strconv.Itoa(int(communityID))+"/leave", nil, strangerToken)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreatePost_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "poster@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "post-community")

	w := performJSONRequest(env.router, http.MethodPost, "/api/posts", map[string]any{
		"title":        "First post",
		"content":      "Hello vitilo",
		"community_id": communityID,
		"post_type":    "text",
	}, token)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreatePost_InvalidCommunity_ReturnsNotFound(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "poster2@ufl.edu")

	w := performJSONRequest(env.router, http.MethodPost, "/api/posts", map[string]any{
		"title":        "Broken post",
		"content":      "missing community",
		"community_id": 9999,
		"post_type":    "text",
	}, token)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetPosts_ReturnsPaginatedFeed(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "feed@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "feed-community")
	_ = createPostWithToken(t, env, token, communityID, "Post 1")
	_ = createPostWithToken(t, env, token, communityID, "Post 2")

	w := performJSONRequest(env.router, http.MethodGet, "/api/posts?page=1&limit=10&sort=new", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateComment_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "commenter@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "comment-community")
	postID := createPostWithToken(t, env, token, communityID, "Comment post")

	w := performJSONRequest(env.router, http.MethodPost, "/api/comments", map[string]any{
		"content": "Nice post",
		"post_id": postID,
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateComment_InvalidParent_ReturnsNotFound(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "commenter2@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "comment-parent-community")
	postID := createPostWithToken(t, env, token, communityID, "Comment parent post")

	w := performJSONRequest(env.router, http.MethodPost, "/api/comments", map[string]any{
		"content":   "Reply",
		"post_id":   postID,
		"parent_id": 9999,
	}, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteComment_ForbiddenForAnotherUser(t *testing.T) {
	env := setupTestEnv(t)
	ownerToken, _ := registerAndGetToken(t, env, "ownercomment@ufl.edu")
	otherToken, _ := registerAndGetToken(t, env, "othercomment@ufl.edu")
	communityID := createCommunityWithToken(t, env, ownerToken, "comment-delete-community")
	postID := createPostWithToken(t, env, ownerToken, communityID, "Delete comment post")
	commentID := createCommentWithToken(t, env, ownerToken, postID, nil, "Owner comment")

	w := performJSONRequest(env.router, http.MethodDelete, "/api/comments/"+strconv.Itoa(int(commentID)), nil, otherToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestVotePost_UpvoteThenToggleOff(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "voter@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "vote-community")
	postID := createPostWithToken(t, env, token, communityID, "Vote post")

	upvote := performJSONRequest(env.router, http.MethodPost, "/api/posts/"+strconv.Itoa(int(postID))+"/vote", map[string]any{"value": 1}, token)
	if upvote.Code != http.StatusOK {
		t.Fatalf("expected first vote 200, got %d body=%s", upvote.Code, upvote.Body.String())
	}

	toggle := performJSONRequest(env.router, http.MethodPost, "/api/posts/"+strconv.Itoa(int(postID))+"/vote", map[string]any{"value": 1}, token)
	if toggle.Code != http.StatusOK {
		t.Fatalf("expected toggle vote 200, got %d body=%s", toggle.Code, toggle.Body.String())
	}
}

func TestVoteComment_Upvote_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "commentvote@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "comment-vote-community")
	postID := createPostWithToken(t, env, token, communityID, "Comment vote post")
	commentID := createCommentWithToken(t, env, token, postID, nil, "Vote this comment")

	w := performJSONRequest(env.router, http.MethodPost, "/api/comments/"+strconv.Itoa(int(commentID))+"/vote", map[string]any{"value": 1}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUploadImage_InvalidExtension_ReturnsBadRequest(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "upload@ufl.edu")

	req := createMultipartUploadRequest(t, "/api/posts/upload", "image", "bad.txt", "not an image", token)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDeletePost_ForbiddenForNonOwner(t *testing.T) {
	env := setupTestEnv(t)
	ownerToken, _ := registerAndGetToken(t, env, "postowner@ufl.edu")
	otherToken, _ := registerAndGetToken(t, env, "postother@ufl.edu")
	communityID := createCommunityWithToken(t, env, ownerToken, "delete-post-community")
	postID := createPostWithToken(t, env, ownerToken, communityID, "Protected post")

	w := performJSONRequest(env.router, http.MethodDelete, "/api/posts/"+strconv.Itoa(int(postID)), nil, otherToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetUserPosts_ReturnsPostsForRequestedUser(t *testing.T) {
	env := setupTestEnv(t)
	token, userID := registerAndGetToken(t, env, "profileposts@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "user-posts-community")
	_ = createPostWithToken(t, env, token, communityID, "Profile Post 1")

	w := performJSONRequest(env.router, http.MethodGet, fmt.Sprintf("/api/users/%d/posts", userID), nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCleanupUploadsDirectoryPath(t *testing.T) {
	// Lightweight sanity test so the single file also touches helper path behavior.
	path := filepath.Join("uploads", "sample.png")
	if filepath.Dir(path) != "uploads" {
		t.Fatalf("expected uploads directory, got %s", filepath.Dir(path))
	}
}

// ==================== NEW TEST CASES ====================

func TestGetAllCommunities_ReturnsOK(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "communities@ufl.edu")
	_ = createCommunityWithToken(t, env, token, "test-community")

	w := performJSONRequest(env.router, http.MethodGet, "/api/communities", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetCommunityByID_ReturnsOK(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "getcommunity@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "community-detail")

	w := performJSONRequest(env.router, http.MethodGet, "/api/communities/"+strconv.Itoa(int(communityID)), nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetCommunityByID_InvalidID_ReturnsNotFound(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "getcommunitybadid@ufl.edu")

	w := performJSONRequest(env.router, http.MethodGet, "/api/communities/9999", nil, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetUsersCommunities_ReturnsOK(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "usercommunities@ufl.edu")
	_ = createCommunityWithToken(t, env, token, "user-comm-1")
	_ = createCommunityWithToken(t, env, token, "user-comm-2")

	w := performJSONRequest(env.router, http.MethodGet, "/api/communities/user/joined", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDownvotePost_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "downvoter@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "downvote-community")
	postID := createPostWithToken(t, env, token, communityID, "Downvote post")

	w := performJSONRequest(env.router, http.MethodPost, "/api/posts/"+strconv.Itoa(int(postID))+"/vote", map[string]any{"value": -1}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDownvotePost_ThenToggleOff(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "downvotetoggle@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "downvote-toggle-community")
	postID := createPostWithToken(t, env, token, communityID, "Toggle downvote post")

	downvote := performJSONRequest(env.router, http.MethodPost, "/api/posts/"+strconv.Itoa(int(postID))+"/vote", map[string]any{"value": -1}, token)
	if downvote.Code != http.StatusOK {
		t.Fatalf("expected first downvote 200, got %d body=%s", downvote.Code, downvote.Body.String())
	}

	toggle := performJSONRequest(env.router, http.MethodPost, "/api/posts/"+strconv.Itoa(int(postID))+"/vote", map[string]any{"value": -1}, token)
	if toggle.Code != http.StatusOK {
		t.Fatalf("expected toggle downvote 200, got %d body=%s", toggle.Code, toggle.Body.String())
	}
}

func TestCreateReplyComment_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "replier@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "reply-community")
	postID := createPostWithToken(t, env, token, communityID, "Reply post")
	parentCommentID := createCommentWithToken(t, env, token, postID, nil, "Parent comment")

	w := performJSONRequest(env.router, http.MethodPost, "/api/comments", map[string]any{
		"content":   "This is a reply",
		"post_id":   postID,
		"parent_id": parentCommentID,
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateReplyComment_DeepNesting_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "deeplyreplier@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "deep-reply-community")
	postID := createPostWithToken(t, env, token, communityID, "Deep reply post")
	level1Comment := createCommentWithToken(t, env, token, postID, nil, "Level 1 comment")
	level2Comment := createCommentWithToken(t, env, token, postID, &level1Comment, "Level 2 comment")

	w := performJSONRequest(env.router, http.MethodPost, "/api/comments", map[string]any{
		"content":   "Level 3 reply",
		"post_id":   postID,
		"parent_id": level2Comment,
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestLeaveCommunity_Success(t *testing.T) {
	env := setupTestEnv(t)
	ownerToken, _ := registerAndGetToken(t, env, "owner3@ufl.edu")
	memberToken, _ := registerAndGetToken(t, env, "member3@ufl.edu")
	communityID := createCommunityWithToken(t, env, ownerToken, "leave-community")

	// Join the community first
	_ = performJSONRequest(env.router, http.MethodPost, "/api/communities/"+strconv.Itoa(int(communityID))+"/join", nil, memberToken)

	// Leave the community
	w := performJSONRequest(env.router, http.MethodPost, "/api/communities/"+strconv.Itoa(int(communityID))+"/leave", nil, memberToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteComment_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "deletecomment@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "delete-comment-community")
	postID := createPostWithToken(t, env, token, communityID, "Delete comment post")
	commentID := createCommentWithToken(t, env, token, postID, nil, "Comment to delete")

	w := performJSONRequest(env.router, http.MethodDelete, "/api/comments/"+strconv.Itoa(int(commentID)), nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDeletePost_Success(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := registerAndGetToken(t, env, "deletepost@ufl.edu")
	communityID := createCommunityWithToken(t, env, token, "delete-post-community")
	postID := createPostWithToken(t, env, token, communityID, "Post to delete")

	w := performJSONRequest(env.router, http.MethodDelete, "/api/posts/"+strconv.Itoa(int(postID)), nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}
