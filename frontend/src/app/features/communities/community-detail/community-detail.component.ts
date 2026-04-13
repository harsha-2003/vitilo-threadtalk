import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterModule } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { PostService } from '../../../core/services/post.service'; 
import { CommunityService } from '../../../core/services/community.service';
import { Community } from '../../../models/community.model';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';

@Component({
  selector: 'app-community-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatDialogModule,
  ],
  templateUrl: './community-detail.component.html',
  styleUrls: ['./community-detail.component.scss'],
})
export class CommunityDetailComponent implements OnInit {
  isLoading = true;
  error: string | null = null;
  posts: any[] = [];
  isLoadingPosts = false;

  communityId!: number;
  community: Community | null = null;

  constructor(
    private route: ActivatedRoute,
    private communityService: CommunityService,
    private postService: PostService,
    private snackBar: MatSnackBar,
  private dialog: MatDialog
  ) {}

  ngOnInit(): void {
  this.communityId = Number(this.route.snapshot.paramMap.get('id'));

  if (!this.communityId || Number.isNaN(this.communityId)) {
    this.error = 'Invalid community id';
    this.isLoading = false;
    return;
  }

  // Load community
  this.communityService.getCommunity(this.communityId).subscribe({
    next: (c) => {
      this.community = c;
      this.isLoading = false;
    },
    error: (err) => {
      console.error(err);
      this.error = 'Failed to load community';
      this.isLoading = false;
    }
  });

  // Load posts (array)
  this.isLoadingPosts = true;
  this.postService.getPostsByCommunity(this.communityId).subscribe({
    next: (res: any) => {
    this.posts = res?.posts ?? [];     // ✅ this is the important fix
    this.isLoadingPosts = false;
  },
    error: (err) => {
      console.error(err);
      this.posts = [];
      this.isLoadingPosts = false;
    }
  });
}
  getCommunityIcon(): string {
    if (!this.community) return 'assets/community-default.png';
    return this.community.icon_url?.trim()
      ? this.community.icon_url
      : 'assets/community-default.png';
  }
  deletePost(post: any): void {
  const ok = confirm(`Delete "${post.title}"?`);
  if (!ok) return;

  this.postService.deletePost(post.id).subscribe({
    next: () => {
      // remove from UI immediately
      this.posts = this.posts.filter(p => p.id !== post.id);
      this.snackBar.open('Post deleted', 'Close', { duration: 2000 });
    },
    error: (err) => {
      console.error(err);
      this.snackBar.open('Failed to delete post', 'Close', { duration: 3000 });
    }
  });
}
}
