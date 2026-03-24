import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule, MatTabChangeEvent } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

import { PostService } from '../../../core/services/post.service';
import { CommunityService } from '../../../core/services/community.service';
import { Community } from '../../../models/community.model';
import { CreatePostRequest } from '../../../models/post.model';

@Component({
  selector: 'app-create-post-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatTabsModule,
    MatCardModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './create-post-dialog.component.html',
  styleUrls: ['./create-post-dialog.component.scss']
})
export class CreatePostDialogComponent implements OnInit {
  postForm: FormGroup;
  communities: Community[] = [];
  imagePreview: string | null = null;
  selectedFile: File | null = null;
  isUploading = false;
  isSubmitting = false;
  currentPostType: 'text' | 'image' | 'link' = 'text';

  constructor(
    private fb: FormBuilder,
    private dialogRef: MatDialogRef<CreatePostDialogComponent>,
    private postService: PostService,
    private communityService: CommunityService,
    private snackBar: MatSnackBar
  ) {
    this.postForm = this.fb.group({
      title: ['', [Validators.required, Validators.maxLength(300)]],
      content: [''],
      community_id: ['', Validators.required],
      post_type: ['text']
    });
  }

  ngOnInit(): void {
    this.loadCommunities();
  }

  loadCommunities(): void {
    this.communityService.getUserCommunities().subscribe({
      next: (communities) => {
        this.communities = communities;
        if (communities.length > 0) {
          this.postForm.patchValue({ community_id: communities[0].id });
        }
      },
      error: (error) => {
        this.snackBar.open('Failed to load communities', 'Close', { duration: 3000 });
      }
    });
  }

  onTabChange(event: MatTabChangeEvent): void {
    const postTypes = ['text', 'image', 'link'];
    this.currentPostType = postTypes[event.index] as 'text' | 'image' | 'link';
    this.postForm.patchValue({ post_type: this.currentPostType });

    // Update content validators based on post type
    const contentControl = this.postForm.get('content');
    if (this.currentPostType === 'link') {
      contentControl?.setValidators([Validators.required, Validators.pattern(/^https?:\/\/.+/)]);
    } else {
      contentControl?.clearValidators();
    }
    contentControl?.updateValueAndValidity();
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      this.handleFile(input.files[0]);
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    
    if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
      this.handleFile(event.dataTransfer.files[0]);
    }
  }

  handleFile(file: File): void {
    // Validate file type
    if (!file.type.startsWith('image/')) {
      this.snackBar.open('Please select an image file', 'Close', { duration: 3000 });
      return;
    }

    // Validate file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
      this.snackBar.open('Image must be less than 5MB', 'Close', { duration: 3000 });
      return;
    }

    this.selectedFile = file;

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      this.imagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  }

  removeImage(event: Event): void {
    event.stopPropagation();
    this.imagePreview = null;
    this.selectedFile = null;
  }

  async onSubmit(): Promise<void> {
    if (this.postForm.invalid) {
      return;
    }

    this.isSubmitting = true;

    try {
      let imageUrl = '';

      // Upload image if it's an image post
      if (this.currentPostType === 'image' && this.selectedFile) {
        this.isUploading = true;
        const uploadResponse = await this.postService.uploadImage(this.selectedFile).toPromise();
        imageUrl = uploadResponse?.image_url || '';
        this.isUploading = false;
      }

      // Create post
      const postData: CreatePostRequest = {
        title: this.postForm.value.title,
        content: this.postForm.value.content,
        community_id: this.postForm.value.community_id,
        post_type: this.currentPostType,
        image_url: imageUrl
      };

      this.postService.createPost(postData).subscribe({
        next: (post) => {
          this.snackBar.open('Post created successfully!', 'Close', { duration: 2000 });
          this.dialogRef.close(post);
        },
        error: (error) => {
          this.snackBar.open('Failed to create post', 'Close', { duration: 3000 });
          this.isSubmitting = false;
        }
      });
    } catch (error) {
      this.snackBar.open('Failed to upload image', 'Close', { duration: 3000 });
      this.isSubmitting = false;
      this.isUploading = false;
    }
  }
}
