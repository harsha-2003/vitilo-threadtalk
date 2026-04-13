# Sprint 3:

Frontend Video - https://youtu.be/D7Wzj7w16eU

Backend Video - https://youtu.be/_bM9zyhTGt4

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
<img width="998" height="948" alt="image" src="https://github.com/user-attachments/assets/1634b47b-a07d-455e-97da-702460cba813" />
<img width="1064" height="998" alt="image" src="https://github.com/user-attachments/assets/88ce44f7-1ee4-4613-8f3b-1545f244eb98" />
<img width="1108" height="1007" alt="image" src="https://github.com/user-attachments/assets/08638b2a-4844-486d-89a4-b460328f8eb1" />


# Sprint 3 Frontend Documentation — Vitilo ThreadTalk

## Overview
Vitilo ThreadTalk is a Reddit-style discussion platform built for a Software Engineering project. The **Sprint 3 Frontend** continues the Angular-based UI developed in earlier sprints and extends it to support richer **community workflows** and **community-scoped posts**, matching the Sprint 3 backend capabilities.

This sprint focused on:
- A complete Communities browsing experience
- A dedicated Community page (`/community/:id`) to display community details and posts
- A dedicated “Create Post” flow inside a community (`/community/:id/create-post`)
- Join/leave community UX integration
- Cypress end-to-end (smoke/unit-style) UI tests for key flows

> Backend reference: Sprint 3 Go API using Gin, GORM, SQLite, JWT auth.

---

## Frontend Technology Stack
- **Framework:** Angular (standalone components)
- **Language:** TypeScript
- **UI Library:** Angular Material
- **Routing:** Angular Router
- **HTTP:** Angular HttpClient
- **State:** Component-level state + service calls (no global store)
- **Testing (E2E / UI tests):** Cypress (`cypress run`)
- **Styling:** SCSS

---

## Sprint 3 Frontend Work Completed

### Communities
Implemented and stabilized the community workflow:
- Browse all communities (`/communities`)
- View communities joined by the current user (tab inside `/communities`)
- Search/filter communities client-side
- Create a new community (Angular Material dialog)
- Join/leave communities from the list UI
- Navigate to community page from community cards (`/community/:id`)

### Community Page (Detail)
Created a dedicated Community page that displays:
- Community header (name, description, icon)
- Member count
- Posts belonging to that community
- “Create Post” CTA linking to `/community/:id/create-post`

### Community Posts
Added community-scoped posting and retrieval:
- Create a post for a specific community
- Fetch posts filtered by `community_id`
- Display community posts under the community page

### Delete Posts (Community Page)
Added a post deletion option to the community posts list:
- Uses `DELETE /api/posts/:id` via `PostService.deletePost()`
- Removes deleted post from UI list without full reload
- Backend enforces ownership validation (only author can delete)

---

## Frontend Route Documentation

### Base URL (Frontend)
By default:
- `http://localhost:4200`

### Main Routes
| Route | Description |
|------|-------------|
| `/` | Landing page |
| `/login` | Login page |
| `/register` | Register page |
| `/feed` | Global feed posts |
| `/post/:id` | Post detail page |
| `/communities` | Communities list page (All + My Communities tabs) |
| `/community/:id` | Community detail page (header + community posts) |
| `/community/:id/create-post` | Create a post specifically inside the community |
| `/profile` | Profile page |

> Important: Community routes must be defined **above** the wildcard `**` route in `app.routes.ts`.

Example community route configuration:

```ts
{
  path: 'community/:id',
  loadComponent: () =>
    import('./features/communities/community-detail/community-detail.component')
      .then(m => m.CommunityDetailComponent)
},
{
  path: 'community/:id/create-post',
  loadComponent: () =>
    import('./features/communities/community-create-post/community-create-post.component')
      .then(m => m.CommunityCreatePostComponent)
},
```

---

## Services Used (Core Services Layer)

All services are organized under:
- `src/app/core/services`

### CommunityService (key methods)
- `getCommunities()`
- `getCommunity(id)`
- `createCommunity(...)`
- `joinCommunity(id)`
- `leaveCommunity(id)`
- `getUserCommunities()`

### PostService (key methods)
Located at:
- `src/app/core/services/post.service.ts`

Key methods used in Sprint 3 frontend:
- `getPosts(page, limit, sort, communityId?)`
- `createPost(request)`
- `deletePost(id)`

#### Fetch posts by community (recommended approach)
Sprint 3 uses the existing pagination endpoint with `community_id` filter:
- `GET /api/posts?community_id=:id`

Frontend leverages:
- `PostService.getPosts(..., communityId)`

---

## Community Posts Integration (End-to-End)

### Creating a Post in a Community
UI flow:
1. User navigates to `/community/:id`
2. Clicks **Create Post**
3. Fills title + content
4. Submits the form
5. Redirects back to `/community/:id`
6. Newly created post appears in the community posts list

**Critical requirement:** the create post payload must include the community id:

```json
{
  "title": "My post title",
  "content": "My post content",
  "community_id": 1
}
```

---

## UI Components Added / Updated (Sprint 3)

### CommunityListComponent
Path:
- `src/app/features/communities/community-list/community-list.component.ts`

Responsibilities:
- Render All Communities + My Communities
- Search/filter
- Join/leave communities
- Open Create Community dialog
- Navigate to `/community/:id` via router

### CommunityDetailComponent
Path:
- `src/app/features/communities/community-detail/community-detail.component.ts`

Responsibilities:
- Load community header data by ID
- Load community posts filtered by ID
- Render posts and “Create Post” CTA
- Support post deletion (frontend) using PostService + backend ownership rules

### CommunityCreatePostComponent
Path:
- `src/app/features/communities/community-create-post/community-create-post.component.ts`

Responsibilities:
- Create a new post in a specific community
- Ensure `community_id` is included
- Redirect to `/community/:id` after success

---

## Cypress UI Tests (Sprint 3)
Cypress tests were extended in Sprint 3 to include Communities page smoke coverage and stable UI assertions.

### Test Runner
From `frontend/`:
```bash
npm run cy:run
```

### Test File
- `cypress/e2e/smoke-login.cy.ts`

### Added Test Block (Sprint 3)
The following tests were added under:

`describe('ThreadTalk - Communities & Community Posts (E2E)', () => { ... })`

> These tests were written to avoid fragile selectors and to assert stable UI behavior.

#### Cypress Test Cases (Sprint 3 Frontend)
```ts
describe('ThreadTalk - Communities & Community Posts (E2E)', () => {
  it('shows the communities page header and create community button', () => {
    cy.visit('/communities');

    cy.get('.communities-page').should('be.visible');
    cy.contains('h1', 'Communities').should('be.visible');
    cy.contains('Discover and join communities at UF').should('be.visible');

    cy.contains('button', 'Create Community').should('be.visible');
  });

  it('shows the search field and allows typing + clearing the query', () => {
    cy.visit('/communities');

    cy.get('input[placeholder="Search by name or description"]')
      .should('be.visible')
      .type('test')
      .should('have.value', 'test');

    // Clear using the suffix close icon button (only appears when searchQuery is non-empty)
    cy.get('button[mat-icon-button][matSuffix]').click();
    cy.get('input[placeholder="Search by name or description"]').should('have.value', '');
  });

  it('renders both tabs: All Communities and My Communities', () => {
    cy.visit('/communities');

    cy.contains('All Communities').should('be.visible');
    cy.contains('My Communities').should('be.visible');
  });

  it('switches to "My Communities" tab and shows empty state OR community cards', () => {
    cy.visit('/communities');

    cy.contains('My Communities').click();

    cy.get('body').then(($body) => {
      const bodyText = $body.text();

      const emptyStateShown =
        bodyText.includes('No communities yet') ||
        bodyText.includes('Join communities to see them here');

      if (emptyStateShown) {
        expect(emptyStateShown).to.eq(true);
      } else {
        cy.get('mat-card.community-card').should('have.length.at.least', 1);
      }
    });
  });

  it('create community dialog opens and can be closed (without submitting)', () => {
    cy.visit('/communities');

    cy.contains('button', 'Create Community').click();
    cy.contains('Create a Community').should('be.visible');

    cy.contains('button', 'Cancel').click();
    cy.contains('Create a Community').should('not.exist');
  });
});
```

---

## Notes / Known Limitations
- Cypress navigation tests that click “View Community” can fail if the UI text differs or if the element is rendered as an `<a>` instead of a `<button>`.  
  Recommendation: add stable selectors like `data-cy="view-community"` to community cards.

- Community member list (actual users) is not rendered if no backend endpoint exists for `/api/communities/:id/members`. The UI displays **member_count**.

- Delete Post is protected by backend ownership validation. If the logged-in user is not the owner, backend should return a `403` and UI should show an error message.

---

## How to Run the Frontend
From `frontend/`:
```bash
npm install
npm start
```

By default:
- Frontend: `http://localhost:4200`

---

## How to Run Cypress Tests
From `frontend/`:
```bash
npm run cy:run
```

---

## Sprint 3 Demo Checklist (Frontend)
During demo/presentation, show:
1. Login / register
2. Navigate to `/communities`
3. Search communities
4. Open and close “Create Community” dialog
5. Join a community
6. Click into community page `/community/:id`
7. Create a post in a community
8. Verify post appears under that community
9. Delete a post (as owner)

---

## Future Improvements
- Add stable `data-cy` selectors for Cypress
- Add community members list UI if backend exposes members endpoint
- Add edit/update post and edit/update comment features
- Add pagination controls on community posts list
- Improve consistency of API response shapes and frontend typing


<img width="1146" height="597" alt="image" src="https://github.com/user-attachments/assets/ea0c99bc-f003-446f-9013-bd316a635003" />
<img width="918" height="717" alt="image" src="https://github.com/user-attachments/assets/6dd566da-9ea8-4bd2-b05d-1f34bf42568d" />
<img width="849" height="687" alt="image" src="https://github.com/user-attachments/assets/e8fcffaa-35fa-49e3-b3f9-6389dc9c583f" />
<img width="903" height="588" alt="image" src="https://github.com/user-attachments/assets/01504f19-428b-49a9-8882-23d771d08511" />





