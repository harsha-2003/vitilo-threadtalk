import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';

import { HeaderComponent } from '../../components/header/header.component';
import { AuthService } from '../../core/services/auth.service';
import { PostService } from '../../core/services/post.service';
import { User } from '../../models/user.model';
import { Post } from '../../models/post.model';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule, HeaderComponent],
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.scss']
})
export class ProfilePageComponent implements OnInit {
  user: User | null = null;

  posts: Post[] = [];
  isLoadingPosts = false;
  error: string | null = null;

  constructor(
    private authService: AuthService,
    private postService: PostService
  ) {}

  ngOnInit(): void {
    this.user = this.authService.getCurrentUser();

    if (!this.user) {
      this.error = 'Please login to view your profile.';
      return;
    }

    // IMPORTANT: use /users/:id/posts
    this.loadUserPosts(this.user.id);
  }

  loadUserPosts(userId: number): void {
    this.isLoadingPosts = true;
    this.error = null;

    this.postService.getUserPosts(userId).subscribe({
      next: (res) => {
        this.posts = res?.posts ?? [];
        this.isLoadingPosts = false;
      },
      error: (err) => {
        console.error('Failed to load user posts', err);
        this.posts = [];
        this.isLoadingPosts = false;
        this.error = 'Failed to load your posts';
      }
    });
  }

  getInitial(): string {
    const s = this.user?.anonymous_username || this.user?.email || '';
    return s ? s.charAt(0).toUpperCase() : '?';
  }
}
