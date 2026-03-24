import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

import { Comment } from '../../../models/comment.model';
import { CommentService } from '../../../core/services/comment.service';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-comment-thread',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatMenuModule,
    MatSnackBarModule
  ],
  templateUrl: './comment-thread.component.html',
  styleUrls: ['./comment-thread.component.scss']
})
export class CommentThreadComponent implements OnInit {
  @Input() comment!: Comment;
  @Input() postId!: number;
  @Input() isNested = false;
  @Output() commentAdded = new EventEmitter<void>();

  replyForm: FormGroup;
  showReplyForm = false;
  isCollapsed = false;

  constructor(
    private fb: FormBuilder,
    private commentService: CommentService,
    private authService: AuthService,
    private snackBar: MatSnackBar
  ) {
    this.replyForm = this.fb.group({
      content: ['', Validators.required]
    });
  }

  ngOnInit(): void {}

  getAvatarUrl(hash: string): string {
    return `https://api.dicebear.com/7.x/identicon/svg?seed=${hash}`;
  }

  getTimeAgo(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return 'just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + 'h ago';
    if (seconds < 2592000) return Math.floor(seconds / 86400) + 'd ago';
    return Math.floor(seconds / 2592000) + 'mo ago';
  }

  onUpvote(): void {
    const newVote = this.comment.user_vote === 1 ? 0 : 1;
    this.vote(newVote);
  }

  onDownvote(): void {
    const newVote = this.comment.user_vote === -1 ? 0 : -1;
    this.vote(newVote);
  }

  private vote(value: number): void {
    const previousVote = this.comment.user_vote;
    const previousCount = this.comment.vote_count;

    // Optimistic update
    this.comment.user_vote = value;
    this.comment.vote_count = previousCount - previousVote + value;

    this.commentService.voteComment(this.comment.id, value).subscribe({
      next: (response) => {
        this.comment.vote_count = response.vote_count;
      },
      error: () => {
        // Revert on error
        this.comment.user_vote = previousVote;
        this.comment.vote_count = previousCount;
        this.snackBar.open('Failed to vote', 'Close', { duration: 2000 });
      }
    });
  }

  toggleReply(): void {
    this.showReplyForm = !this.showReplyForm;
    if (!this.showReplyForm) {
      this.replyForm.reset();
    }
  }

  toggleCollapse(): void {
    this.isCollapsed = !this.isCollapsed;
  }

  submitReply(): void {
    if (this.replyForm.invalid) return;

    const replyData = {
      content: this.replyForm.value.content,
      post_id: this.postId,
      parent_id: this.comment.id
    };

    this.commentService.createComment(replyData).subscribe({
      next: (reply) => {
        if (!this.comment.replies) {
          this.comment.replies = [];
        }
        this.comment.replies.push(reply);
        this.replyForm.reset();
        this.showReplyForm = false;
        this.snackBar.open('Reply added!', 'Close', { duration: 2000 });
        this.commentAdded.emit();
      },
      error: () => {
        this.snackBar.open('Failed to add reply', 'Close', { duration: 3000 });
      }
    });
  }

  onReplyAdded(): void {
    this.commentAdded.emit();
  }

  isOwnComment(): boolean {
    const currentUser = this.authService.getCurrentUser();
    return currentUser?.id === this.comment.user_id;
  }

  deleteComment(): void {
    if (confirm('Are you sure you want to delete this comment?')) {
      this.commentService.deleteComment(this.comment.id).subscribe({
        next: () => {
          this.snackBar.open('Comment deleted', 'Close', { duration: 2000 });
          this.commentAdded.emit();
        },
        error: () => {
          this.snackBar.open('Failed to delete comment', 'Close', { duration: 3000 });
        }
      });
    }
  }

  reportComment(): void {
    this.snackBar.open('Comment reported. Thank you!', 'Close', { duration: 3000 });
  }
}
