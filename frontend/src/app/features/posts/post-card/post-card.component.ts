import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';

import { Post } from '../../../models/post.model';
import { PostService } from '../../../core/services/post.service';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-post-card',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatSnackBarModule
  ],
  templateUrl: './post-card.component.html',
  styleUrls: ['./post-card.component.scss']
})
export class PostCardComponent {
  @Input() post!: Post;
  @Output() voteChange = new EventEmitter<{ postId: number; voteCount: number }>();
  
  isExpanded = false;

  constructor(
    private postService: PostService,
    private authService: AuthService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  onUpvote(): void {
    const newVote = this.post.user_vote === 1 ? 0 : 1;
    this.vote(newVote);
  }

  onDownvote(): void {
    const newVote = this.post.user_vote === -1 ? 0 : -1;
    this.vote(newVote);
  }

  private vote(value: number): void {
    const previousVote = this.post.user_vote;
    const previousCount = this.post.vote_count;

    // Optimistic update
    this.post.user_vote = value;
    this.post.vote_count = previousCount - previousVote + value;

    this.postService.votePost(this.post.id, value).subscribe({
      next: (response) => {
        this.post.vote_count = response.vote_count;
        this.voteChange.emit({ postId: this.post.id, voteCount: response.vote_count });
      },
      error: (error) => {
        // Revert on error
        this.post.user_vote = previousVote;
        this.post.vote_count = previousCount;
        this.snackBar.open('Failed to vote. Please try again.', 'Close', { duration: 3000 });
      }
    });
  }

  formatVoteCount(count: number): string {
    if (count >= 10000) {
      return (count / 1000).toFixed(1) + 'k';
    } else if (count >= 1000) {
      return (count / 1000).toFixed(1) + 'k';
    }
    return count.toString();
  }

  getTimeAgo(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return 'just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + ' minutes ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + ' hours ago';
    if (seconds < 2592000) return Math.floor(seconds / 86400) + ' days ago';
    if (seconds < 31536000) return Math.floor(seconds / 2592000) + ' months ago';
    return Math.floor(seconds / 31536000) + ' years ago';
  }

  toggleExpand(): void {
    this.isExpanded = !this.isExpanded;
  }

  openImage(imageUrl: string): void {
    window.open(imageUrl, '_blank');
  }

  sharePost(): void {
    const url = `${window.location.origin}/post/${this.post.id}`;
    if (navigator.clipboard) {
      navigator.clipboard.writeText(url);
      this.snackBar.open('Link copied to clipboard!', 'Close', { duration: 2000 });
    }
  }

  savePost(): void {
    this.snackBar.open('Post saved!', 'Close', { duration: 2000 });
  }

  reportPost(): void {
    this.snackBar.open('Post reported. Thank you for your feedback.', 'Close', { duration: 3000 });
  }

  isOwnPost(): boolean {
    const currentUser = this.authService.getCurrentUser();
    return currentUser?.id === this.post.user_id;
  }

  editPost(): void {
    this.snackBar.open('Edit feature coming soon!', 'Close', { duration: 2000 });
  }

  deletePost(): void {
    if (confirm('Are you sure you want to delete this post?')) {
      this.postService.deletePost(this.post.id).subscribe({
        next: () => {
          this.snackBar.open('Post deleted successfully', 'Close', { duration: 2000 });
          // Emit event or refresh feed
        },
        error: (error) => {
          this.snackBar.open('Failed to delete post', 'Close', { duration: 3000 });
        }
      });
    }
  }
}
