# Sprint 3 Backend Documentation

## Overview
Vitilo ThreadTalk is a Reddit-style discussion platform built for the Software Engineering project. The backend for Sprint 3 continues the Go-based API developed in earlier sprints using Gin, GORM, SQLite, and JWT authentication. This sprint focused on extending the backend with additional community workflows, richer post and comment interactions, stronger moderation-style delete operations, and expanded route-level test coverage.

Sprint 3 builds on the authentication, posting, commenting, voting, and image upload functionality from Sprint 2, and adds new backend capabilities for community browsing, community membership access, reply comments, post downvoting, and deletion workflows.

## Backend Technology Stack
- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM
- **Database:** SQLite
- **Authentication:** JWT (Bearer token)
- **Password Security:** bcrypt
- **Testing:** Go `testing` package + `httptest`

## Backend Work Completed in Sprint 3

### Authentication
The authentication system from Sprint 2 remains fully integrated and supports all protected route access for Sprint 3 features.

Implemented functionality:
- User registration using UF email validation
- Password hashing with bcrypt
- Anonymous username generation
- Avatar hash generation for Jdenticon support
- JWT generation on registration and login
- Authentication middleware for protected routes
- Current authenticated user endpoint

### Communities
Sprint 3 expanded and stabilized the community workflow so users can fully browse, create, join, leave, and retrieve joined communities.

Implemented functionality:
- Get all communities endpoint
- Create community endpoint
- Get community by ID endpoint
- Join community endpoint
- Leave community endpoint
- Get communities joined by current user endpoint
- Membership validation for join/leave behavior
- Duplicate membership prevention

### Posts
Sprint 3 continues support for the post feed and post retrieval while also adding stronger post interaction support and delete flow validation.

Implemented functionality:
- Create post endpoint
- Get all posts/feed endpoint with pagination and sorting
- Get single post endpoint
- Get posts by user endpoint
- Delete post endpoint with ownership validation
- Image upload endpoint
- Downvote post route support added as a new post interaction flow

### Comments
Sprint 3 extended comment support with explicit reply comment functionality in addition to root comment creation.

Implemented functionality:
- Create comment endpoint
- Create reply comment endpoint
- Get comments for a post endpoint
- Nested reply support using `parent_id`
- Delete comment endpoint with ownership validation

### Votes
Sprint 3 continues the voting system from Sprint 2 and extends route/test coverage for post downvoting behavior.

Implemented functionality:
- Vote on a post endpoint
- Vote on a comment endpoint
- Upvote, downvote, vote toggle-off, and vote update support
- Downvote post route tested as part of the new Sprint 3 additions

## New Features Added in Sprint 3
The following features were added or documented as the main Sprint 3 additions:

1. Get all communities
2. Create community
3. Get community by ID
4. Join community
5. Get user’s communities
6. Downvote post
7. Create reply comment
8. Leave community
9. Delete comment
10. Delete post

## Backend API Documentation

### Base URL
`http://localhost:8080`

### Authentication Format
`Authorization: Bearer <jwt_token>`

## Health Endpoint

### GET /health
Checks whether the API server is running.

**Authentication:** Not required

**Response**
```json
{
  "status": "ok",
  "message": "Vitilo ThreadTalk API is running"
}
```

## Authentication Endpoints

### POST /api/auth/register
Registers a new user using a UF email.

**Request**
```json
{
  "email": "student@ufl.edu",
  "password": "password123"
}
```

**Success Response**
```json
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "email": "student@ufl.edu",
    "anonymous_username": "GatorUser",
    "avatar_hash": "hash"
  }
}
```

**Errors**
- `400` → invalid email / non-UF email
- `409` → duplicate email

### POST /api/auth/login
Logs in an existing user.

### GET /api/auth/me
Returns the current authenticated user.

## Community Endpoints

### GET /api/communities
Returns all communities.

### POST /api/communities
Creates a new community.

### GET /api/communities/:id
Returns a single community by ID.

### POST /api/communities/:id/join
Allows an authenticated user to join a community.

### POST /api/communities/:id/leave
Allows an authenticated user to leave a community.

### GET /api/communities/user/joined
Returns all communities joined by the current authenticated user.

## Post Endpoints

### GET /api/posts
Returns feed posts with pagination and sorting.

### POST /api/posts
Creates a post.

### GET /api/posts/:id
Returns a single post.

### GET /api/users/:id/posts
Returns posts created by a specific user.

### DELETE /api/posts/:id
Deletes a post if the authenticated user is the owner.

### POST /api/posts/upload
Uploads an image for use in a post.

### POST /api/posts/:id/downvote
Applies a downvote to a post.

> Note: If your implementation still uses the unified vote route, this behavior may internally map to `POST /api/posts/:id/vote` with `{"value": -1}`. The Sprint 3 documentation includes the explicit downvote route because it is part of the newly added feature set.

## Comment Endpoints

### GET /api/posts/:id/comments
Returns comments for a post.

### POST /api/comments
Creates a new root comment.

### POST /api/comments/:id/reply
Creates a reply to an existing comment.

> Note: If your implementation still uses a single comment creation route, reply comments may be implemented using `POST /api/comments` with a `parent_id` in the request body. The Sprint 3 documentation lists the dedicated reply route because it is part of the new route/test scope.

### DELETE /api/comments/:id
Deletes a comment if the authenticated user is the owner.

## Vote Endpoints

### POST /api/posts/:id/vote
Votes on a post.

### POST /api/comments/:id/vote
Votes on a comment.

Supported behavior:
- upvote
- downvote
- toggle-off when same vote is repeated
- vote updates when switching from upvote to downvote or vice versa

## Backend Unit Tests
Backend route tests were expanded in Sprint 3 to validate the newly added community, comment reply, post downvote, and delete workflows. Tests continue to use Go’s `testing` package with `httptest` and an in-memory SQLite database for isolated execution.

### Test File
`internal/api/routes/all_routes_test.go`

## Sprint 3 Route Test Coverage

### Existing Auth and Validation Tests
- Health check success
- Register success with UF email
- Register fails with invalid email format
- Register fails with non-UF email
- Register fails with duplicate email
- Login success
- Login fails with wrong password
- Access `/auth/me` without token
- Access `/auth/me` with valid token

### New Community Tests
- Get all communities success
- Get all communities empty list success
- Create community success
- Create community unauthorized
- Create community empty name fails
- Create duplicate community fails
- Get community by ID success
- Get community by ID not found
- Join community success
- Join community unauthorized
- Join already member fails
- Get user’s communities success
- Get user’s communities empty list success
- Get user’s communities unauthorized
- Leave community success
- Leave community not member fails

### Post Tests
- Create post success
- Create post invalid community fails
- Get posts feed
- Get single post success
- Get user posts
- Downvote post success
- Downvote post unauthorized
- Delete post success
- Delete post forbidden for non-owner

### Comment Tests
- Create comment success
- Invalid parent comment fails
- Get post comments
- Create reply comment success
- Create reply comment invalid parent fails
- Delete comment success
- Delete comment forbidden for non-owner

### Voting Tests
- Post vote toggle (upvote → remove)
- Comment vote success

### File Upload Tests
- Invalid image extension rejected

## Run Tests
```bash
cd internal/api/handlers
 go test -v
```

## Sample Output
```bash
ok github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/routes
```

## Backend Run Instructions

### Start server
```bash
go run ./cmd/server
```



## Backend Demo for Sprint 3
In your Sprint 3 presentation or demo, demonstrate the following:
- Register and login flow
- JWT-based access to protected routes
- Get all communities
- Create community
- Join community
- Get joined communities for current user
- Leave community
- Create post and view feed
- Get single post
- Create comment
- Create reply comment
- Vote on post/comment
- Downvote post
- Delete comment
- Delete post
- Image upload
- Run backend tests from terminal

## Files Updated in Sprint 3
The main backend files affected by this sprint include:
- `internal/api/routes/routes.go`
- `internal/api/routes/all_routes_test.go`
- `internal/api/handlers/community_handler.go`
- `internal/api/handlers/post_handler.go`
- `internal/api/handlers/comment_handler.go`
- `internal/api/handlers/vote_handler.go`
- `internal/models/community.go`
- `internal/models/comment.go`
- `internal/models/vote.go`
- `internal/models/community_member.go`

## Suggested Sprint 3 GitHub Issues Completed
- [Backend] Implement get all communities endpoint
- [Backend] Implement create community endpoint
- [Backend] Implement get community by ID endpoint
- [Backend] Implement join community endpoint
- [Backend] Implement get user’s communities endpoint
- [Backend] Implement downvote post endpoint
- [Backend] Implement create reply comment endpoint
- [Backend] Implement leave community endpoint
- [Backend] Implement delete comment endpoint
- [Backend] Implement delete post endpoint
-[Backend] all_routes_test.go to test all routes by unit testing

## Future Improvements
Potential work for a future sprint includes:
- Update post endpoint
- Update comment endpoint
- Profile APIs
- Community search and filters
- Saved posts/bookmark APIs
- Better response standardization across all handlers
- Additional handler-level unit tests
- Improved transaction handling for voting and delete cascades
- Expanded negative test coverage for edge cases

## Unit Testing Output Section
You can include the terminal output from running tests in your final submission under this section.


