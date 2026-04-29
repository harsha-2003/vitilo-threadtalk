# Sprint 4 Frontend Documentation — Vitilo ThreadTalk
Youtube  - https://youtu.be/7Z5YiT2zGwQ
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
