# Vitilo ThreadTalk Backend API Documentation

Base URL: `http://localhost:8080`

The backend is a Go API built with Gin, GORM, SQLite, JWT authentication, and CORS support. All `/api` routes except registration and login require an `Authorization: Bearer <token>` header.

## Health

### `GET /health`
Returns `200 OK` when the API is running.

```json
{
  "status": "ok",
  "message": "Vitilo ThreadTalk API is running"
}
```

## Authentication

### `POST /api/auth/register`
Creates a user account. The email must use the configured UF email domain.

Request:

```json
{
  "email": "student@ufl.edu",
  "password": "password123"
}
```

Success `201 Created`:

```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "email": "student@ufl.edu",
    "anonymous_username": "HappyGator42",
    "avatar_hash": "32-character-avatar-hash"
  }
}
```

Errors: `400` invalid input or non-UF email, `409` duplicate email, `500` server error.

### `POST /api/auth/login`
Authenticates an existing user.

Request:

```json
{
  "email": "student@ufl.edu",
  "password": "password123"
}
```

Success `200 OK`: same response shape as registration.

Errors: `400` invalid input, `401` invalid credentials.

### `GET /api/auth/me`
Returns the authenticated user.

Success `200 OK`:

```json
{
  "id": 1,
  "email": "student@ufl.edu",
  "anonymous_username": "HappyGator42",
  "avatar_hash": "32-character-avatar-hash"
}
```

## Communities

### `GET /api/communities`
Returns all communities with membership metadata for the current user.

### `POST /api/communities`
Creates a community and automatically joins the creator.

Request:

```json
{
  "name": "uf-tech",
  "description": "UF tech club",
  "icon_url": "https://example.com/icon.png"
}
```

Success `201 Created`:

```json
{
  "id": 1,
  "name": "uf-tech",
  "description": "UF tech club",
  "icon_url": "https://example.com/icon.png",
  "member_count": 1,
  "is_member": true,
  "created_at": "2026-04-28T15:00:00Z"
}
```

### `GET /api/communities/user/joined`
Returns communities joined by the authenticated user.

### `GET /api/communities/:id`
Returns one community by id.

### `POST /api/communities/:id/join`
Joins a community. Returns `409 Conflict` if the user is already a member.

### `POST /api/communities/:id/leave`
Leaves a community. Returns `404 Not Found` if the user is not a member.

## Posts

### `GET /api/posts`
Returns a paginated feed.

Query parameters:

| Name | Default | Description |
| --- | --- | --- |
| `page` | `1` | Page number. |
| `limit` | `20` | Page size, capped at 100. |
| `sort` | `new` | `new`, `top`, or `hot`. |
| `community_id` | none | Optional community filter. |

Success `200 OK`:

```json
{
  "posts": [],
  "total": 0,
  "page": 1,
  "limit": 20,
  "total_pages": 0
}
```

### `POST /api/posts`
Creates a text or image post.

Request:

```json
{
  "title": "First post",
  "content": "Hello vitilo",
  "community_id": 1,
  "post_type": "text",
  "image_url": "/uploads/example.png"
}
```

Success `201 Created` returns:

```json
{
  "id": 1,
  "title": "First post",
  "content": "Hello vitilo",
  "image_url": "",
  "post_type": "text",
  "vote_count": 0,
  "comment_count": 0,
  "created_at": "2026-04-28T15:00:00Z",
  "user_id": 1,
  "anonymous_username": "HappyGator42",
  "avatar_hash": "32-character-avatar-hash",
  "community_id": 1,
  "community_name": "uf-tech",
  "user_vote": 0
}
```

### `POST /api/posts/upload`
Uploads an image using multipart form data with field name `image`.

Allowed extensions: `.jpg`, `.jpeg`, `.png`, `.gif`. Maximum size: 5 MB.

Success `200 OK`:

```json
{
  "image_url": "/uploads/generated-file-name.png",
  "message": "Image uploaded successfully"
}
```

### `GET /api/posts/:id`
Returns a single post.

### `DELETE /api/posts/:id`
Deletes a post owned by the authenticated user. Non-owners receive `403 Forbidden`.

### `GET /api/users/:id/posts`
Returns posts created by a specific user.

## Comments

### `GET /api/posts/:id/comments`
Returns top-level comments for a post with nested replies.

### `POST /api/comments`
Creates a comment or reply.

Request:

```json
{
  "content": "Nice post",
  "post_id": 1,
  "parent_id": null
}
```

Success `201 Created`:

```json
{
  "id": 1,
  "content": "Nice post",
  "vote_count": 0,
  "created_at": "2026-04-28T15:00:00Z",
  "user_id": 1,
  "anonymous_username": "HappyGator42",
  "avatar_hash": "32-character-avatar-hash",
  "post_id": 1,
  "parent_id": null,
  "user_vote": 0
}
```

### `DELETE /api/comments/:id`
Deletes a comment owned by the authenticated user. Non-owners receive `403 Forbidden`.

## Votes

### `POST /api/posts/:id/vote`
Upvotes, downvotes, changes a vote, or toggles the same vote off for a post.

Request:

```json
{
  "value": 1
}
```

`value` must be `1` or `-1`.

Success `200 OK`:

```json
{
  "message": "Vote processed successfully",
  "vote_count": 1
}
```

### `POST /api/comments/:id/vote`
Matches post voting behavior for comments.

## Static Files

### `GET /uploads/:filename`
Serves uploaded post images from the local `uploads` directory.

## Backend Test Command

```powershell
.\run_tests.ps1
```

Equivalent raw command:

```powershell
$env:GOCACHE = Join-Path (Get-Location) ".gocache"
go test -v ./... -timeout 60s
```
