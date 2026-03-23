# Sprint 2 – Vitilo ThreadTalk Backend

## Overview
Vitilo ThreadTalk is a Reddit-style discussion platform built for the Software Engineering project. The backend for Sprint 2 was implemented in Go using Gin, GORM, SQLite, and JWT authentication. This sprint focused on delivering the core backend APIs required for user authentication, communities, posts, comments, voting, and image upload.

---

## Backend Technology Stack
- **Language:** Go  
- **Framework:** Gin  
- **ORM:** GORM  
- **Database:** SQLite  
- **Authentication:** JWT (Bearer token)  
- **Password Security:** bcrypt  
- **Testing:** Go `testing` package + `httptest`  

---

## Backend Work Completed in Sprint 2

### Authentication
- User registration using UF email validation  
- Password hashing with bcrypt  
- Anonymous username generation  
- Avatar hash generation for Jdenticon support  
- JWT token generation on register and login  
- Protected route access using authentication middleware  
- Current authenticated user endpoint  

### Communities
- Community model and membership model  
- Create community endpoint  
- Get all communities endpoint  
- Get community by ID endpoint  
- Join community endpoint  
- Leave community endpoint  
- Get communities joined by current user endpoint  

### Posts
- Post model  
- Create post endpoint  
- Get feed posts endpoint with pagination and sorting  
- Get single post endpoint  
- Get posts by user endpoint  
- Delete post endpoint with ownership validation  
- Image upload endpoint  

### Comments
- Comment model with nested reply support using `parent_id`  
- Create comment endpoint  
- Get comments for a post endpoint  
- Delete comment endpoint with ownership validation  

### Votes
- Vote model  
- Vote on a post endpoint  
- Vote on a comment endpoint  
- Support for upvote, downvote, toggle-off, and vote updates  

---

# Backend API Documentation

## Base URL
http://localhost:8080

## Authentication Format
Authorization: Bearer <jwt_token>

---

## Health Endpoint

### GET /health
Checks whether the API server is running.

**Authentication:** Not required  

**Response**
{
  "status": "ok",
  "message": "Vitilo ThreadTalk API is running"
}

---

## Authentication Endpoints

### POST /api/auth/register
Registers a new user using UF email.

Request:
{
  "email": "student@ufl.edu",
  "password": "password123"
}

Success:
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "email": "student@ufl.edu",
    "anonymous_username": "GatorUser",
    "avatar_hash": "hash"
  }
}

Errors:
- 400 → invalid email / not UF domain  
- 409 → duplicate email  

---

### POST /api/auth/login
Login existing user.

---

### GET /api/auth/me
Returns current user.

---

## Community Endpoints

### GET /api/communities
Get all communities.

### POST /api/communities
Create new community.

### POST /api/communities/:id/join
Join community.

### POST /api/communities/:id/leave
Leave community.

---

## Post Endpoints

### GET /api/posts
Get posts feed with pagination.

### POST /api/posts
Create post.

### GET /api/posts/:id
Get single post.

### DELETE /api/posts/:id
Delete post (owner only).

### POST /api/posts/upload
Upload image.

---

## Comment Endpoints

### GET /api/posts/:id/comments
Get comments.

### POST /api/comments
Create comment.

### DELETE /api/comments/:id
Delete comment.

---

## Vote Endpoints

### POST /api/posts/:id/vote
Vote on post.

### POST /api/comments/:id/vote
Vote on comment.

---

# Backend Unit Tests

The backend unit tests were implemented using Go’s `testing` package and `httptest` to simulate HTTP requests. An in-memory SQLite database is used to ensure isolated test execution.

All test cases are consolidated into a single file for simplicity.

---

## Test File
internal/api/routes/all_routes_test.go

---

## Test Coverage (20+ Tests)

### Authentication & Validation
- Health check success  
- Register success with UF email  
- Register fails with invalid email format  
- Register fails with non-UF email  
- Register fails with duplicate email  
- Login success  
- Login fails with wrong password  
- Access `/auth/me` without token  
- Access `/auth/me` with valid token  

---

### Community
- Create community success  
- Create duplicate community fails  
- Join community success  
- Join already member fails  
- Leave community not member fails  

---

### Posts
- Create post success  
- Create post invalid community fails  
- Get posts feed  
- Delete post forbidden (non-owner)  

---

### Comments
- Create comment success  
- Invalid parent comment fails  
- Delete comment forbidden  

---

### Voting
- Post vote toggle (upvote → remove)  
- Comment vote success  

---

### File Upload
- Invalid image extension rejected  

---

### Users
- Get user posts  

---

## Run Tests
go test ./...

### Sample Output
ok github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/routes

---

# Backend Run Instructions

### Start server
go run ./cmd/server

### Run tests
go test ./...

---

# Backend Demo (Video Requirements)

In your presentation, demonstrate:
- Register + login flow  
- JWT authentication  
- Community creation and join/leave  
- Post creation and feed  
- Comment creation  
- Voting system  
- Image upload  
- Running backend tests  

---

# Future Improvements

- Update post endpoint  
- Update comment endpoint  
- Profile APIs  
- Community search & filters  
- Transactional voting improvements  
- Better error handling standardization  
- More test coverage  
