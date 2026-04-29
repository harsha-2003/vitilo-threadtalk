import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatCardModule } from '@angular/material/card';

import { PostService } from '../../../core/services/post.service';
import { CommentService } from '../../../core/services/comment.service';
import { AuthService } from '../../../core/services/auth.service';
import { Post } from '../../../models/post.model';
import { Comment } from '../../../models/comment.model';
import { User } from '../../../models/user.model';

@Component({
  selector: 'app-post-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonToggleModule,
    MatSnackBarModule,
    MatCardModule
  ],
  templateUrl: './post-detail.component.html',
  styleUrls: ['./post-detail.component.scss']
})
export class PostDetailComponent implements OnInit {
  post: Post | null = null;
  comments: Comment[] = [];
  currentUser: User | null = null;
  commentForm: FormGroup;
  replyForm: FormGroup;
  isLoading = true;
  isSubmitting = false;
  error: string | null = null;
  commentSort = 'best';
  replyingTo: number | null = null;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private postService: PostService,
    private commentService: CommentService,
    private authService: AuthService,
    private snackBar: MatSnackBar
  ) {
    this.commentForm = this.fb.group({
      content: ['', [Validators.required, Validators.minLength(1)]]
    });
    
    this.replyForm = this.fb.group({
      content: ['', [Validators.required, Validators.minLength(1)]]
    });
  }

  ngOnInit(): void {
    console.log('🚀 Post Detail Component Initialized');
    this.currentUser = this.authService.getCurrentUser();
    console.log('👤 Current User:', this.currentUser);
    
    const postId = this.route.snapshot.paramMap.get('id');
    console.log('📝 Post ID from route:', postId);
    
    if (postId) {
      this.loadPost(+postId);
      this.loadComments(+postId);
    }
  }

  loadPost(postId: number): void {
    console.log(`📥 Loading post ${postId}...`);
    this.isLoading = true;
    this.postService.getPost(postId).subscribe({
      next: (post) => {
        console.log('✅ Post loaded:', post);
        this.post = post;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('❌ Error loading post:', error);
        this.error = 'Failed to load post';
        this.isLoading = false;
      }
    });
  }

  loadComments(postId: number): void {
    console.log(`📥 Loading comments for post ${postId}...`);
    this.commentService.getComments(postId).subscribe({
      next: (comments) => {
        console.log('✅ Comments loaded:', comments);
        console.log('📊 Number of comments:', comments.length);
        this.comments = this.sortComments(comments);
        console.log('📊 Sorted comments:', this.comments);
      },
      error: (error) => {
        console.error('❌ Failed to load comments:', error);
        this.comments = [];
      }
    });
  }

  // MISSING METHOD - Add Comment
  submitComment(): void {
    if (this.commentForm.invalid || !this.post) {
      console.warn('⚠️ Comment form invalid or post missing');
      return;
    }

    this.isSubmitting = true;

    const commentData = {
      content: this.commentForm.value.content.trim(),
      post_id: this.post.id,
      parent_id: null
    };

    console.log('📤 Submitting comment:', commentData);

    this.commentService.createComment(commentData).subscribe({
      next: (comment) => {
        console.log('✅ Comment created:', comment);
        this.comments.unshift(comment);
        this.commentForm.reset();
        this.isSubmitting = false;
        this.snackBar.open('Comment added! 💬', 'Close', { duration: 2000 });
        
        if (this.post) {
          this.post.comment_count++;
        }
      },
      error: (error) => {
        console.error('❌ Failed to add comment:', error);
        this.isSubmitting = false;
        const errorMsg = error.error?.error || error.message || 'Unknown error';
        this.snackBar.open(`Failed to add comment: ${errorMsg}`, 'Close', { duration: 3000 });
      }
    });
  }

  startReply(commentId: number): void {
    console.log('💬 Starting reply to comment:', commentId);
    this.replyingTo = commentId;
    this.replyForm.reset();
  }

  cancelReply(): void {
    console.log('❌ Cancelling reply');
    this.replyingTo = null;
    this.replyForm.reset();
  }

  submitReply(parentId: number): void {
    if (this.replyForm.invalid || !this.post) {
      return;
    }

    const replyData = {
      content: this.replyForm.value.content.trim(),
      post_id: this.post.id,
      parent_id: parentId
    };

    console.log('📤 Submitting reply:', replyData);

    this.commentService.createComment(replyData).subscribe({
      next: (reply) => {
        console.log('✅ Reply created:', reply);
        const parent = this.findComment(this.comments, parentId);
        if (parent) {
          if (!parent.replies) {
            parent.replies = [];
          }
          parent.replies.push(reply);
        }
        this.cancelReply();
        this.snackBar.open('Reply added!', 'Close', { duration: 2000 });
        
        if (this.post) {
          this.post.comment_count++;
        }
      },
      error: (error) => {
        console.error('❌ Failed to add reply:', error);
        this.snackBar.open('Failed to add reply', 'Close', { duration: 3000 });
      }
    });
  }

  deleteComment(commentId: number): void {
    if (!confirm('Delete this comment?')) {
      return;
    }

    console.log('🗑️ Deleting comment:', commentId);

    this.commentService.deleteComment(commentId).subscribe({
      next: () => {
        console.log('✅ Comment deleted:', commentId);
        this.comments = this.removeComment(this.comments, commentId);
        if (this.post) {
          this.post.comment_count--;
        }
        this.snackBar.open('Comment deleted', 'Close', { duration: 2000 });
      },
      error: (error) => {
        console.error('❌ Failed to delete comment:', error);
        this.snackBar.open('Failed to delete comment', 'Close', { duration: 3000 });
      }
    });
  }

  voteComment(commentId: number, value: number): void {
    console.log(`🗳️ Voting on comment ${commentId} with value ${value}`);
    
    this.commentService.voteComment(commentId, value).subscribe({
      next: (response: any) => {
        console.log('✅ Vote recorded:', response);
        const comment = this.findComment(this.comments, commentId);
        if (comment) {
          comment.vote_count = response.vote_count;
          comment.user_vote = value;
        }
      },
      error: (error) => {
        console.error('❌ Failed to vote:', error);
      }
    });
  }

  votePost(value: number): void {
    if (!this.post) return;

    console.log(`🗳️ Voting on post ${this.post.id} with value ${value}`);

    this.postService.votePost(this.post.id, value).subscribe({
      next: (response: any) => {
        console.log('✅ Post vote recorded:', response);
        if (this.post) {
          this.post.vote_count = response.vote_count;
          this.post.user_vote = value;
        }
      },
      error: (error) => {
        console.error('❌ Failed to vote on post:', error);
      }
    });
  }

  onSortChange(event: any): void {
    this.commentSort = event.value;
    console.log('🔄 Sorting comments by:', this.commentSort);
    this.comments = this.sortComments(this.comments);
  }

  sortComments(comments: Comment[]): Comment[] {
    if (!comments || comments.length === 0) {
      console.log('⚠️ No comments to sort');
      return [];
    }

    console.log(`🔄 Sorting ${comments.length} comments by ${this.commentSort}`);

    switch (this.commentSort) {
      case 'new':
        return [...comments].sort((a, b) => 
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
      case 'top':
        return [...comments].sort((a, b) => b.vote_count - a.vote_count);
      case 'best':
      default:
        return [...comments].sort((a, b) => {
          const scoreA = a.vote_count / (Math.log10(Date.now() - new Date(a.created_at).getTime()) + 1);
          const scoreB = b.vote_count / (Math.log10(Date.now() - new Date(b.created_at).getTime()) + 1);
          return scoreB - scoreA;
        });
    }
  }

  private findComment(comments: Comment[], id: number): Comment | null {
    for (const comment of comments) {
      if (comment.id === id) {
        return comment;
      }
      if (comment.replies) {
        const found = this.findComment(comment.replies, id);
        if (found) return found;
      }
    }
    return null;
  }

  private removeComment(comments: Comment[], id: number): Comment[] {
    return comments.filter(comment => {
      if (comment.id === id) {
        return false;
      }
      if (comment.replies) {
        comment.replies = this.removeComment(comment.replies, id);
      }
      return true;
    });
  }

  getTimeAgo(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return 'just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + 'h ago';
    return Math.floor(seconds / 86400) + 'd ago';
  }

  formatVotes(votes: number): string {
    if (!votes) return '0';
    if (votes >= 1000) {
      return (votes / 1000).toFixed(1) + 'k';
    }
    return votes.toString();
  }

  goBack(): void {
    this.router.navigate(['/feed']);
  }
}
