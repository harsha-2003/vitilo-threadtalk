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
 # output of unit testing
 <img width="1245" height="874" alt="Screenshot 2026-03-25 001122" src="https://github.com/user-attachments/assets/6970df45-3a92-4ec4-b226-ce1e97836326" />
 <img width="1016" height="870" alt="Screenshot 2026-03-25 001142" src="https://github.com/user-attachments/assets/15a01ea4-1a29-4592-86b5-7e0afbfc4386" />
<img width="1126" height="885" alt="Screenshot 2026-03-25 001159" src="https://github.com/user-attachments/assets/c5330890-2a89-4a97-a504-43cfd05ff9d5" />
 <img width="1140" height="476" alt="Screenshot 2026-03-25 001230" src="https://github.com/user-attachments/assets/3c49fbf9-4490-494e-bcdd-1711b498c794" />


# Sprint 2 – Vitilo ThreadTalk Frontend

## Overview
Vitilo ThreadTalk is a Reddit-style discussion platform built for the Software Engineering project. The **frontend for Sprint 2** was implemented using **Angular (Standalone Components)** with **Angular Material** styling. This sprint focused on delivering the core user-facing UI flows and integrating them with backend APIs, along with adding both **unit tests (Angular/Jasmine/Karma)** and a **Cypress E2E test suite** for the Login page.

---

## Frontend Technology Stack
- **Language:** TypeScript  
- **Framework:** Angular (Standalone Components)  
- **Styling/UI:** SCSS + Angular Material  
- **Routing:** Angular Router  
- **Auth Storage:** LocalStorage (JWT token + user object)  
- **Testing (Unit):** Jasmine + Karma (via Angular CLI)  
- **Testing (E2E):** Cypress  

---

## Frontend Work Completed in Sprint 2

### Authentication UI (Login/Register)
- Built/validated Login page UI using Angular Reactive Forms  
- Email and password validation:
  - Email required + email format validation
  - Password required + minlength validation (8 characters)
- Loading state and disabled submit behavior based on form validity  
- Password visibility toggle (hide/show password)
- Navigation links:
  - “Create Account” → `/register`
  - “Back to Home” → `/`

### Header & Navigation
- Profile navigation support from header (profile icon/button routes to Profile page)
- Fixed compilation issues caused by duplicate method definitions in header (`goToProfile`) by ensuring a single function implementation

### Profile Page
- Profile page UI created/updated to display:
  - User information (anonymous username, avatar initial/hash)
  - User’s posts section with loading + empty state
- Integrated profile posts API using a user-id route pattern:
  - `GET /api/users/:id/posts`
- Updated PostService to fetch posts for a specific user based on logged-in user id

### General Integration / Build Fixes
- Resolved Angular build errors due to:
  - TemplateUrl missing files (NG2008)
  - Routes syntax issues (TS1005)
  - Template variable mismatches (`isLoading` vs `isLoadingPosts`)
- Standardized component state variables and ensured template bindings match component properties

---

# Frontend E2E Test Documentation (Cypress)

## Base URL
http://localhost:4200

## Cypress Suite (Sprint 2)
**Test File:**  
`frontend/cypress/e2e/smoke-login.cy.ts`

### What the Cypress suite validates
The Cypress tests focus on **UI correctness** and **form behavior** without requiring a successful backend login response.

**Key scenarios tested (Login Page):**
- Login page renders correctly (title, subtitle, structure)
- Email input behavior:
  - placeholder, type, autocomplete, clearing input
  - invalid email shows validation message
- Password input behavior:
  - placeholder, autocomplete, minlength validation
  - required validation error
- Submit button behavior:
  - disabled while form is invalid
  - enabled only when form is valid
  - has correct CSS classes and attributes
- Password visibility toggle:
  - exists and switches icon state
  - toggles visibility twice returns to original state
- Navigation:
  - “Create Account” navigates to `/register`
  - “Back to Home” navigates to `/`

---

# Frontend Unit Tests (Angular)

Unit tests were written using Angular’s testing utilities with Jasmine/Karma. The goal was to cover frontend logic at the **function level** and to follow the guideline:  
> **Aim for a 1:1 unit test-to-function ratio** where possible.

## Unit Test Targets (examples used in Sprint 2)
- **AuthService**
  - `isAuthenticated()`
  - `getCurrentUser()`
  - `logout()`
- **PostService**
  - `getUserPosts(userId)`

> Note: Unit tests are structured to validate logic deterministically (URL construction, local storage behavior, return values) without relying on real backend state.

---

# Run Instructions (Frontend)

## Start Frontend
```bash
cd frontend
npm install
ng serve
```

Frontend runs at:
- http://localhost:4200

## Run Unit Tests
```bash
cd frontend
ng test --watch=false
```

## Run Cypress Tests
1) Start Angular:
```bash
cd frontend
ng serve
```

2) Run Cypress headless:
```bash
cd frontend
npm run cy:run
```

---

# Demo (Video Requirements – Frontend Portion)

In the narrated presentation, demonstrate:

### UI Functionality
- Login page:
  - form validations (invalid email, short password)
  - password visibility toggle
  - navigation to Register and Home
- Profile page:
  - profile button navigation from header
  - user information display
  - posts section with empty state or populated posts (depending on seeded data)

### Test Results
- Show unit test results:
  - `ng test --watch=false`
- Show Cypress results:
  - `npm run cy:run`

### Suggested Split for Team Narration
- **Member 1:** UI walkthrough + integration points (Login/Profile)
- **Member 2:** Testing walkthrough (unit tests + Cypress run output)

---

# Future Improvements (Frontend)
- Add unit tests for additional components (community list, feed, post detail, comments)
- Add Cypress E2E tests for:
  - Register flow
  - Create post flow
  - Comment flow
  - Vote interactions
- Improve test stability by adding dedicated `data-cy` selectors for key UI elements
- Add profile edit flow UI and tests
