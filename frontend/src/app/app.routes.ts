import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./features/auth/landing/landing.component').then(m => m.LandingComponent)
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./features/auth/register/register.component').then(m => m.RegisterComponent)
  },
  {
    path: 'feed',
    loadComponent: () =>
      import('./features/posts/post-feed/post-feed.component').then(m => m.PostFeedComponent)
  },
  {
    path: 'post/:id',
    loadComponent: () =>
      import('./features/posts/post-detail/post-detail.component').then(m => m.PostDetailComponent)
  },
  {
    path: 'communities',
    loadComponent: () =>
      import('./features/communities/community-list/community-list.component').then(m => m.CommunityListComponent)
  },
  {
    path: 'profile',
    loadComponent: () =>
      import('./features/profile/profile-page.component').then(m => m.ProfilePageComponent)
  },
  
 {
  path: 'community/:id',
  loadComponent: () =>
    import('./features/communities/community-detail/community-detail.component')
      .then(m => m.CommunityDetailComponent)
},
{
  path: 'community/:id/create-post',
  loadComponent: () =>
    import('../app/features/communities/community-create-post/community-create-post.component')
      .then(m => m.CommunityCreatePostComponent)
},
{
    path: '**',
    redirectTo: ''
  },
];
