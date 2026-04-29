# Sprint 4 Frontend Documentation — Vitilo ThreadTalk
Youtube  - https://youtu.be/7Z5YiT2zGwQ
Youtube - https://youtu.be/f6D0DWbjxMA
# Sprint 4 Backend Documentation

## Overview

Vitilo ThreadTalk is a Reddit-style anonymous discussion platform built for the Software Engineering project. The backend for Sprint 4 continues the Go-based API developed in earlier sprints using Gin, GORM, SQLite, JWT authentication, and bcrypt password security.

Sprint 4 focused on preparing the backend for final submission by improving backend documentation, validating handler-level backend behavior, updating test execution scripts, and creating a clear backend API reference for the presentation and grading submission.

This sprint builds on the authentication, community, post, comment, vote, image upload, and delete workflows completed in Sprint 3. The main Sprint 4 backend work was documentation, verification, testing cleanup, and final project readiness.

## Backend Technology Stack

- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM
- **Database:** SQLite
- **Authentication:** JWT Bearer token
- **Password Security:** bcrypt
- **Testing:** Go `testing` package + `httptest`
- **API Documentation:** Markdown documentation in `API_DOCS.md`

## Backend Work Completed in Sprint 4

### Documentation

Sprint 4 added and updated documentation needed for final submission.

Completed work:

- Created updated backend API documentation in `API_DOCS.md`
- Added endpoint descriptions for auth, communities, posts, comments, votes, uploads, and health checks
- Documented request and response examples for the main API workflows
- Documented JWT authentication format for protected routes
- Added backend test command instructions
- Updated this Sprint 4 report using the Sprint 3 documentation style

### Backend Test Execution

Sprint 4 improved how backend tests are executed so the submission video can show the handler test suite clearly.

Completed work:

- Updated `run_tests.ps1` to run handler tests with `go test -v ./internal/api/handlers -timeout 60s`
- Updated `run_tests_Version3.sh` to run all backend packages with `go test -v ./... -timeout 60s`
- Configured both test scripts to use a local `.gocache` directory
- Added `.gocache/` to `.gitignore`
- Verified the handler test suite from the backend project root

### Handler Test Reliability

Sprint 4 cleaned up backend handler tests so they run reliably through the PowerShell test command.

Completed work:

- Updated handler tests to use isolated in-memory SQLite database names
- Configured the PowerShell test runner to run only the handler test package
- Confirmed that `internal/api/handlers/all_handlers_test.go` passes through `run_tests.ps1`

### Final Backend Review

Sprint 4 included a full backend review for the final presentation.

Completed work:

- Reviewed implemented backend routes
- Reviewed authentication flow
- Reviewed protected route behavior
- Reviewed ownership checks for deleting posts and comments
- Reviewed vote toggle and downvote behavior
- Reviewed image upload validation
- Prepared backend video presentation notes

## Backend Functionality Completed by Final Sprint

### Authentication

Implemented functionality:

- User registration using UF email validation
- Password hashing with bcrypt
- Anonymous username generation
- Avatar hash generation for anonymous profile icons
- JWT generation on registration and login
- Authentication middleware for protected routes
- Current authenticated user endpoint

### Communities

Implemented functionality:

- Get all communities endpoint
- Create community endpoint
- Get community by ID endpoint
- Join community endpoint
- Leave community endpoint
- Get communities joined by current user endpoint
- Membership validation for join/leave behavior
- Duplicate membership prevention
- Community response fields for member count and current-user membership status

### Posts

Implemented functionality:

- Create post endpoint
- Get all posts/feed endpoint
- Feed pagination
- Feed sorting with `new`, `top`, and `hot`
- Optional community filtering
- Get single post endpoint
- Get posts by user endpoint
- Delete post endpoint with ownership validation
- Image upload endpoint
- Post response fields for vote count, comment count, user vote, author anonymity, and community metadata

### Comments

Implemented functionality:

- Create comment endpoint
- Create reply comment using `parent_id`
- Get comments for a post endpoint
- Nested reply support
- Delete comment endpoint with ownership validation
- Comment response fields for vote count, user vote, author anonymity, and replies

### Votes

Implemented functionality:

- Vote on a post endpoint
- Vote on a comment endpoint
- Upvote support
- Downvote support
- Vote toggle-off when the same vote is repeated
- Vote update when switching from upvote to downvote or downvote to upvote
- Updated aggregate vote counts

### Image Upload

Implemented functionality:

- Multipart image upload endpoint
- Upload field name: `image`
- Accepted extensions: `.jpg`, `.jpeg`, `.png`, `.gif`
- Maximum file size: 5 MB
- Uploaded file serving through `/uploads`

## New Sprint 4 Work

The following items were completed or documented as the main Sprint 4 additions:

1. Updated Sprint 4 backend documentation
2. Created complete backend API documentation
3. Documented JWT authentication behavior
4. Documented backend request and response formats
5. Updated PowerShell backend test runner
6. Updated shell backend test runner
7. Added local Go cache support for tests
8. Ignored local Go cache files in git
9. Improved handler test database isolation
10. Configured the PowerShell test runner for handler tests
11. Verified the handler test suite
12. Documented the handler test command
13. Listed backend handler unit tests
14. Prepared handler test output for the final submission
15. Prepared backend video presentation notes
16. Added final GitHub issue list for Sprint 4 tracking

## Backend API Documentation

### Base URL

`http://localhost:8080`

### Authentication Format

`Authorization: Bearer <jwt_token>`

All `/api` routes require authentication except:

- `POST /api/auth/register`
- `POST /api/auth/login`

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

Registers a new user using a UF email address.

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
    "anonymous_username": "HappyGator42",
    "avatar_hash": "hash"
  }
}
```

**Errors**

- `400` invalid email, short password, or non-UF email
- `409` duplicate email
- `500` server error

### POST /api/auth/login

Logs in an existing user and returns a JWT token.

**Request**

```json
{
  "email": "student@ufl.edu",
  "password": "password123"
}
```

**Errors**

- `400` invalid request body
- `401` invalid credentials

### GET /api/auth/me

Returns the current authenticated user.

**Authentication:** Required

## Community Endpoints

### GET /api/communities

Returns all communities with `member_count` and `is_member` values for the current user.

### POST /api/communities

Creates a new community and automatically joins the creator.

**Request**

```json
{
  "name": "uf-tech",
  "description": "UF tech club",
  "icon_url": "https://example.com/icon.png"
}
```

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

Returns feed posts with pagination, sorting, and optional community filtering.

Supported query parameters:

- `page`
- `limit`
- `sort`
- `community_id`

Supported sort values:

- `new`
- `top`
- `hot`

### POST /api/posts

Creates a post.

**Request**

```json
{
  "title": "First post",
  "content": "Hello vitilo",
  "community_id": 1,
  "post_type": "text",
  "image_url": "/uploads/example.png"
}
```

### GET /api/posts/:id

Returns a single post.

### GET /api/users/:id/posts

Returns posts created by a specific user.

### DELETE /api/posts/:id

Deletes a post if the authenticated user is the owner.

### POST /api/posts/upload

Uploads an image for use in a post.

**Multipart field name:** `image`

**Allowed file types:**

- `.jpg`
- `.jpeg`
- `.png`
- `.gif`

**Maximum file size:** 5 MB

### POST /api/posts/:id/vote

Votes on a post.

**Request**

```json
{
  "value": 1
}
```

Use `1` for upvote and `-1` for downvote.

## Comment Endpoints

### GET /api/posts/:id/comments

Returns comments for a post.

### POST /api/comments

Creates a new root comment or reply comment.

**Request for root comment**

```json
{
  "content": "Nice post",
  "post_id": 1
}
```

**Request for reply comment**

```json
{
  "content": "This is a reply",
  "post_id": 1,
  "parent_id": 2
}
```

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
- toggle-off when the same vote is repeated
- vote updates when switching from upvote to downvote or downvote to upvote

## Backend Unit Tests

Backend tests validate handler-level logic for authentication, communities, posts, comments, votes, and avatar hash generation. Tests use Go's `testing` package, `httptest`, and in-memory SQLite databases.

### Handler Test File

`internal/api/handlers/all_handlers_test.go`

## Sprint 4 Backend Test Coverage

### Handler Tests

- Register success
- Register invalid email
- Register non-UF email
- Register short password
- Register duplicate email
- Login success
- Login wrong password
- Login user not found
- Get current user success
- Get current user not found
- Create community
- Get communities
- Join community
- Leave community
- Get community
- Create post
- Create post with empty title fails
- Create post with invalid community fails
- Get posts
- Get single post
- Delete post
- Delete post forbidden for non-owner
- Get user posts
- Get my posts
- Create comment
- Create comment with invalid post fails
- Get comments
- Delete comment
- Delete comment forbidden for non-owner
- Vote post upvote
- Vote post downvote
- Vote post invalid value fails
- Vote comment
- Vote comment invalid comment fails
- Generate avatar hash

## Run Tests

### PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File .\run_tests.ps1
```

### Raw Go Command

```powershell
$env:GOCACHE = Join-Path (Get-Location) ".gocache"
go test -v ./internal/api/handlers -timeout 60s
```

## Sample Output

```bash
ok github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/handlers
All backend tests succeeded.
```

## Backend Run Instructions

### Start Server

```bash
go run ./cmd/server
```

The server starts on `PORT` from the environment, or `8080` by default.


## Files Updated in Sprint 4

The main backend files affected by this sprint include:

- `Sprint4.md`
- `API_DOCS.md`
- `run_tests.ps1`
- `run_tests_Version3.sh`
- `.gitignore`
- `internal/api/handlers/all_handlers_test.go`

## Suggested Sprint 4 GitHub Issues Completed

1.[Backend] Update Sprint 4 backend documentation (Sprint4.md)
2.[Backend] Create complete backend API documentation
3.[Backend] Document JWT bearer-token authentication requirements
4.[Backend] Add request and response examples for backend API routes
5.[Backend] Implement centralized API response and error handling
6.[Backend] Improve JWT middleware validation and protected routes
7.[Backend] Fix handler test database isolation using test DB setup
8.[Backend] Add backend unit tests for authentication routes
9.[Backend] Add backend unit tests for community APIs
10.[Backend] Add backend unit tests for post APIs
11.[Backend] Add backend unit tests for comment APIs
12.[Backend] Add backend unit tests for vote and upload APIs
13.[Backend] Update PowerShell test runner for handler tests
14.[Backend] Update shell script to run all backend test packages
15.[Backend] Configure local Go cache for faster test execution
16.[Backend] Ignore local Go cache and build artifacts in .gitignore



## Unit Testing Output Section



<img width="1920" height="1200" alt="image" src="https://github.com/user-attachments/assets/87328f90-4a20-4a7e-be9b-2fb6c1f0a430" />
<img width="1920" height="1200" alt="image" src="https://github.com/user-attachments/assets/a2e6ac03-c665-46e6-86e3-808f4cf6315f" />
<img width="1920" height="1200" alt="image" src="https://github.com/user-attachments/assets/be0fd9ab-c9ff-4bb4-8336-f1915019b217" />




## Overview
Vitilo ThreadTalk is a Reddit-style discussion platform built for a Software Engineering project. The **Sprint 4 Frontend** continues the Angular-based Single Page Application (SPA) built in Sprints 1–3 and focuses on:

- Completing and stabilizing remaining Sprint 3 community workflows
- Improving UI reliability and routing stability
- Expanding Cypress end-to-end test coverage (Login + Communities module)
- Preparing final documentation required for Sprint 4 submission (run instructions + test instructions)

Sprint 4 builds on features delivered in Sprint 3 including:
- Communities browsing, join/leave, and community detail page
- Community-scoped post creation and visibility
- Post delete workflows (ownership enforcement handled by backend)

---

## Frontend Technology Stack
- **Framework:** Angular (Standalone Components)
- **Language:** TypeScript
- **UI Library:** Angular Material
- **Routing:** Angular Router
- **HTTP:** Angular HttpClient + Interceptors (JWT injection)
- **State Approach:** Component-level state + service calls (no global store)
- **Styling:** SCSS + dark theme customization
- **Testing (E2E/UI):** Cypress

---

## Sprint 4 Work Completed (Frontend)

### 1) Communities UI Stability Improvements
- Verified communities page renders reliably (`/communities`)
- Stable tab switching between:
  - **All Communities**
  - **My Communities**
- Search input behavior tested (typing does not break routing)
- Community create dialog tested for:
  - required fields visibility
  - cancel behavior
  - submit disabled state when invalid (if validation exists)

### 2) Community Detail Page Stability
- Confirmed deep-link refresh stability for community detail page
- Verified critical sections render after reload:
  - **Members**
  - **Posts**

> Note: Displaying a list of actual member users depends on whether the backend exposes a `GET /api/communities/:id/members` endpoint. If not available, the UI still displays the community `member_count` and renders a placeholder/empty state.

### 3) Expanded Cypress E2E Test Coverage
Sprint 4 includes:
- A full Login Page smoke/validation suite
- Expanded Communities feature tests for UI stability and dialog behavior

---

## Frontend Routes (Key Pages)
| Route | Description |
|------|-------------|
| `/` | Landing page |
| `/login` | Login page |
| `/register` | Register page |
| `/feed` | Global feed |
| `/post/:id` | Post detail page |
| `/communities` | Communities list page (All + My tabs) |
| `/community/:id` | Community detail page (Members + Posts) |
| `/community/:id/create-post` | Create post inside a community |
| `/profile` | Profile page |

---

## Services Used (Core Services Layer)
All frontend services are located in:
- `src/app/core/services`

### CommunityService (Sprint 3–4 usage)
- `getCommunities()`
- `getCommunity(id)`
- `createCommunity(...)`
- `joinCommunity(id)`
- `leaveCommunity(id)`
- `getUserCommunities()`

### PostService (Sprint 3–4 usage)
- `getPosts(page, limit, sort, communityId?)`
- `createPost(request)`
- `deletePost(id)`

**Community posts filter:**
- Uses backend feed endpoint with community filter:
  - `GET /api/posts?community_id=:id`

---

## Cypress UI / E2E Tests (Sprint 4)

### Test Runner
From the `frontend/` directory:
```bash
npm run cy:run
```

### Cypress Test File
- `cypress/e2e/smoke-login.cy.ts`

### Test Suites Included

#### A) Login Page Smoke & Validation Tests
These tests validate:
- Page renders correctly
- Form validation works (invalid email, short password, required fields)
- Password visibility toggle works
- Navigation works:
  - Login → Register
  - Login → Home
- Key Material UI elements render and styling classes exist

#### B) Communities UI Stability + Dialog Tests (Updated Sprint 4 Cases)
These tests validate:
- Communities page loads reliably and main container renders
- Search input behavior does not trigger navigation away
- Tab switching (All ↔ My) does not break UI
- Create Community dialog:
  - opens successfully
  - shows required fields (Community Name, Description)
  - closes via Cancel
  - submit disabled when invalid (if validation exists)
- Navigate away and return to communities page still works
- Community detail deep-link refresh stability (reload does not crash)
  - verifies **Members** and **Posts** headings on the community detail page

**Updated Sprint 4 Cypress tests included (Communities block):**
- `communities page loads without showing a console error (basic health check)`
- `search clear button is not visible before typing (if rendered conditionally)`
- `typing in search does not navigate away from /communities`
- `switching tabs does not break the page`
- `Create Community dialog shows required fields`
- `Create Community submit button disabled when invalid (if supported)`
- `navigating away and back keeps communities page functional`
- `community page supports deep link refresh (reload does not crash)`

---

## How to Run the Frontend (Sprint 4)

### 1) Install dependencies
From `frontend/`:
```bash
npm install
```

### 2) Run the frontend dev server
```bash
ng serve
```

Default:
- Frontend runs at: `http://localhost:4200`

> Ensure the backend is running (default `http://localhost:8080`) so login/community data can load.

---

## How to Run Cypress Tests (Sprint 4)

### 1) Start frontend first
In one terminal:
```bash
cd frontend
npm start
```

### 2) Run Cypress in another terminal
```bash
cd frontend
npm run cy:run
```

If Cypress requires a specific `baseUrl`, confirm your `cypress.config.ts` includes:
- `baseUrl: "http://localhost:4200"`

---

## Notes / Known Limitations
- Some Cypress tests depend on the UI structure (button labels, dialog title). For long-term stability, it is recommended to add `data-cy` selectors such as:
  - `data-cy="create-community-btn"`
  - `data-cy="create-community-dialog"`
  - `data-cy="search-communities-input"`
  - `data-cy="community-tab-all"`
  - `data-cy="community-tab-mine"`
- Displaying an actual member list requires backend endpoint support:
  - `GET /api/communities/:id/members`

---

## Sprint 4 Demo Checklist (Frontend)
During the narrated demo, demonstrate:
1. Login / register pages
2. Navigate to `/communities`
3. Search communities
4. Switch All/My tabs
5. Open and close Create Community dialog
6. View a community page `/community/:id`
7. Create a post inside the community
8. Verify post appears under that community
9. Delete a post (as owner)
10. Show Cypress test run output (`npm run cy:run`)
