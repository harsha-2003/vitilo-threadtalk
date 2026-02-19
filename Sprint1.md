## Demo Videos
- [Backend Demo](https://youtu.be/qWu2sUgLWv4)
- [Frontend Demo](https://youtu.be/YntjwFpkOJA)
- [Github]-(https://github.com/harsha-2003/vitilo-threadtalk)
# Vitilo ThreadTalk - Frontend

## Sprint 1 Report

------------------------------------------------------------------------


## 👥 Team Overview

**Project:** Vitilo ThreadTalk  
**Sprint:** Sprint 1  
**Focus:** Frontend Core Infrastructure & UI Components

Sprint 1 focused on establishing the Angular frontend foundation, building reusable UI components, implementing navigation and layout structure, and creating authentication pages.

------------------------------------------------------------------------

## 🧩 User Stories

------------------------------------------------------------------------

## 🏗️ Epic 1: Frontend Infrastructure & Setup

### 🟢 User Story 1 --- Angular Project Initialization

**As a** developer  
**I want to** initialize the Angular frontend project with standalone components  
**So that** the project follows modern Angular best practices and has a solid foundation

**Acceptance Criteria:**
- Angular project configured with standalone components
- Base application routing set up
- Environment configuration files defined for development and production
- Application builds and runs successfully without errors

**Status:** ✅ Completed

------------------------------------------------------------------------

## 🎨 Epic 2: Layout & Navigation Components

### 🟢 User Story 2 --- Navigation Bar Component

**As a** user  
**I want to** see a navigation bar at the top of the application  
**So that** I can easily access different sections of the platform

**Acceptance Criteria:**
- Responsive navigation bar component created
- Navigation links properly configured
- User profile/avatar display included
- Mobile-friendly design implemented
- Subscription handling and memory leak prevention

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 3 --- Left Sidebar Component

**As a** user  
**I want to** see a sidebar with my subscribed communities and navigation options  
**So that** I can quickly access my favorite communities

**Acceptance Criteria:**
- Reusable and responsive sidebar component created
- Display list of user's subscribed communities
- Support for nested menus
- Collapsible/expandable functionality
- Community creation feature
- Proper subscription management

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 4 --- Right Sidebar Component

**As a** user  
**I want to** see contextual information and tools in a right sidebar  
**So that** I can access additional features relevant to the current page

**Acceptance Criteria:**
- Right sidebar component displaying contextual information
- Modular and collapsible design
- Easy to extend with widgets or dynamic content
- Proper state management
- Memory leak prevention with subscription cleanup

**Status:** ✅ Completed

------------------------------------------------------------------------

## 🔐 Epic 3: Authentication & Landing Pages

### 🟢 User Story 5 --- Landing Page

**As a** new visitor  
**I want to** see an attractive landing page  
**So that** I understand what the platform offers and can navigate to login or register

**Acceptance Criteria:**
- Responsive and user-friendly landing page created
- Clear introduction to the platform (purpose, features, tagline)
- Prominent Call-to-Action (CTA) buttons for Login and Register
- Responsive across desktop and mobile devices
- Clean, reusable UI components
- Proper routing integration

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 6 --- Register Page

**As a** new user  
**I want to** create an account on the platform  
**So that** I can access the community features

**Acceptance Criteria:**
- Standalone Register component created
- Form accepts valid user input and blocks invalid submissions
- Email validation (UF email requirement)
- Password validation
- Error/validation messages display appropriately
- Navigation to Login page works correctly
- No console or routing errors
- Register page loads correctly at `/register`

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 7 --- Login Page

**As a** registered user  
**I want to** log in to my account  
**So that** I can access my personalized content

**Acceptance Criteria:**
- Standalone Login component created
- Form validation for email and password
- Error handling for invalid credentials
- Navigation to Register page available
- Proper routing configuration
- Login page loads correctly at `/login`
- Form submission triggers authentication flow

**Status:** ✅ Completed

------------------------------------------------------------------------

## ✅ Sprint 1 Summary

### Successfully Completed All Issues

All planned frontend issues were successfully completed and merged. The frontend now includes:

#### ✨ Key Achievements:

1. **Angular Project Setup**
   - Standalone components architecture
   - Routing configuration
   - Environment files for dev/prod
   - Build and serve configurations

2. **Core Layout Components**
   - Responsive Navigation Bar with user profile
   - Left Sidebar with community navigation
   - Right Sidebar for contextual information
   - Mobile-friendly responsive design

3. **Code Quality Improvements**
   - Proper subscription management to prevent memory leaks
   - Removal of mock data in favor of service integration
   - Enhanced error handling
   - Improved state management patterns
   - TypeScript best practices

4. **Authentication Pages**
   - Professional Landing Page with CTA buttons
   - Complete Register Page with validation
   - Complete Login Page with form handling
   - Proper routing between auth pages

5. **Angular Best Practices**
   - Standalone components throughout
   - Reactive forms for user input
   - Service-based architecture
   - Proper component lifecycle management
   - Clean code structure and organization

------------------------------------------------------------------------

## 📋 Planned Frontend Issues

Below are the issues planned and completed during Sprint 1:

### Infrastructure & Setup
1. **#51** - [Frontend] Initialize Angular project & routing
2. **#53** - [Frontend] Updated Angular project & routing
3. **#60** - [Frontend] Updated Angular project & routing

### Navigation & Layout Components
4. **#66** - [Frontend] Create Right Sidebar Component
5. **#67** - [Frontend] Create Sidebar Component
6. **#70** - [Frontend] Update Navigation Bar Component
7. **#72** - [Frontend] Update Right Sidebar Component
8. **#78** - [Frontend] Update Sidebar Components

### Component Improvements & Refactoring
9. **#81** - [Frontend] Improve Navbar Component – Subscription Handling, Mock Data Removal & Code Optimization
10. **#82** - [Frontend] Improve Right Sidebar Component – Subscription Handling & State Management
11. **#85** - [Frontend] Refactor Sidebar Component – Subscription Management & Feature Completion

### Landing & Authentication Pages
12. **#76** - [Frontend] Create Landing Page
13. **#89** - [Frontend] Update Landing Page
14. **#90** - [Frontend] Landing Page Component
15. **#94** - [Frontend] Update Register Page
16. **#96** - [Frontend] Create a standalone Register Component
17. **#99** - [Frontend] Update Login Page
18. **#101** - [Frontend] Validate Login Component

------------------------------------------------------------------------

## 🎯 Next Steps for Sprint 2

With the frontend foundation complete, Sprint 2 will focus on:

1. **Post & Community Features**
   - Post creation and display components
   - Community page layouts
   - Voting system UI
   - Comment functionality

2. **State Management**
   - Implement centralized state management (NgRx or signals)
   - User profile state
   - Community subscription state

3. **Enhanced Features**
   - Image upload UI
   - Post sorting and filtering
   - Search functionality
   - Notifications

------------------------------------------------------------------------

## 🏆 Conclusion

Sprint 1 successfully established a robust Angular frontend with a clean component architecture, responsive design, and proper Angular best practices. The team demonstrated excellent collaboration and code quality, setting a strong foundation for future sprints.

All planned issues were completed and merged, representing 100% completion rate for the sprint.



---------------------------------------------------------------------------------------------------------------------------------------------------------------------------

# Vitilo ThreadTalk

## Sprint 1 Report

------------------------------------------------------------------------

## 👥 Team Overview

Project: **Vitilo ThreadTalk**\
Sprint: **Sprint 1**\
Focus: **Backend Core Functionality**

Sprint 1 focused on building the foundational backend architecture
including authentication, community management, post creation, voting,
and database configuration.

------------------------------------------------------------------------

# 🧩 User Stories

------------------------------------------------------------------------

## 🔐 Epic 1: User Authentication & Authorization

### 🟢 User Story 1 --- User Registration

**As a** UF student\
**I want to** register using my `@ufl.edu` email\
**So that** only verified university students can access the platform

**Acceptance Criteria:** - Email must end with `@ufl.edu` - Password
securely hashed using bcrypt - Anonymous username generated - Avatar
hash generated - JWT token issued upon successful registration

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 2 --- User Login

**As a** registered user\
**I want to** log in securely\
**So that** I can access my account

**Acceptance Criteria:** - Credentials validated - JWT token returned -
401 returned on invalid login

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 3 --- JWT-Based Authorization

**As a** logged-in user\
**I want** protected endpoints\
**So that** only authenticated users can perform actions

**Acceptance Criteria:** - Middleware validates JWT - User ID extracted
from token - Protected routes secured

**Status:** ✅ Completed

------------------------------------------------------------------------

# 🌐 Epic 2: Community Management

### 🟢 User Story 4 --- View Communities

**As a** user\
**I want to** see available communities\
**So that** I can browse topics

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 5 --- Create Community

**As a** user\
**I want to** create a community\
**So that** I can start new discussions

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 6 --- Join & Leave Community

**As a** user\
**I want to** join or leave communities\
**So that** I can personalize my feed

**Status:** ✅ Completed

------------------------------------------------------------------------

# 📝 Epic 3: Post Management

### 🟢 User Story 7 --- Create Post

**As a** user\
**I want to** create text or image posts\
**So that** I can participate in discussions

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 8 --- View Posts Feed

**As a** user\
**I want to** browse posts with pagination and sorting\
**So that** I can see discussions efficiently

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 9 --- Delete Own Post

**As a** user\
**I want to** delete my own post\
**So that** I can remove content I created

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 10 --- Upload Image

**As a** user\
**I want to** upload images for posts\
**So that** I can attach media

**Status:** ✅ Completed

------------------------------------------------------------------------

# 🗳️ Epic 4: Voting System

### 🟢 User Story 11 --- Vote on Posts

**As a** user\
**I want to** upvote or remove my vote\
**So that** I can express my opinion

**Status:** ✅ Completed

------------------------------------------------------------------------

### 🟢 User Story 12 --- Vote Model Implementation

**As a** system\
**I need** a vote model\
**So that** voting relationships are stored properly

**Status:** ✅ Completed

------------------------------------------------------------------------

# 🛠️ Infrastructure Stories

### 🟢 Database Configuration & Migration

-   SQLite and GORM configured
-   Auto migrations implemented
-   Tables created successfully

**Status:** ✅ Completed

------------------------------------------------------------------------

# 📌 Planned Backend Issues

1.  Initialize backend structure\
2.  Configure SQLite & GORM\
3.  Add database migrations\
4.  Create User model\
5.  Add UF email validation\
6.  Implement password hashing\
7.  Generate anonymous username\
8.  Generate avatar hash\
9.  Implement JWT token generation\
10. Implement login endpoint\
11. Create JWT middleware\
12. Protect API routes\
13. Create Community model\
14. Implement get communities endpoint\
15. Implement create community endpoint\
16. Implement join/leave endpoint\
17. Create Post model\
18. Implement create post endpoint\
19. Implement get posts endpoint\
20. Implement delete post endpoint\
21. Implement image upload\
22. Create Vote model

------------------------------------------------------------------------

# ✅ Sprint 1 Summary

All planned backend issues were successfully completed and merged via
Pull Requests into the dev branch.

Sprint 1 established a fully functional backend ready for frontend
integration in Sprint 2.
