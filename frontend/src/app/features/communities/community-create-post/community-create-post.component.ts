import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { PostService } from '../../../core/services/post.service';

@Component({
  selector: 'app-community-create-post',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,

    MatCardModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatIconModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: './community-create-post.component.html',
  styleUrls: ['./community-create-post.component.scss'],
})
export class CommunityCreatePostComponent implements OnInit {
  communityId!: number;
  isSubmitting = false;

  form: FormGroup;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private postService: PostService,
    private snackBar: MatSnackBar
  ) {
    // IMPORTANT: initialize form here (after fb is injected)
    this.form = this.fb.group({
      title: ['', [Validators.required, Validators.maxLength(150)]],
      content: ['', [Validators.required, Validators.maxLength(5000)]],
    });
  }

  ngOnInit(): void {
    this.communityId = Number(this.route.snapshot.paramMap.get('id'));

    if (!this.communityId || Number.isNaN(this.communityId)) {
      this.snackBar.open('Invalid community id', 'Close', { duration: 3000 });
      this.router.navigate(['/communities']);
    }
  }

  submit(): void {
    if (this.form.invalid || this.isSubmitting) return;

    this.isSubmitting = true;

    // Payload MUST match backend. Most common is { title, content, community_id }
    const payload: any = {
      title: this.form.value.title,
      content: this.form.value.content,
      community_id: this.communityId,
    };

    this.postService.createPost(payload).subscribe({
      next: () => {
        this.snackBar.open('Post created!', 'Close', { duration: 2500 });
        this.router.navigate(['/community', this.communityId]);
      },
      error: (err) => {
        console.error(err);
        this.isSubmitting = false;
        this.snackBar.open('Failed to create post', 'Close', { duration: 3500 });
      }
    });
  }
}
