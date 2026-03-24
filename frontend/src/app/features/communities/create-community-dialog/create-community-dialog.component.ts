import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

import { CommunityService } from '../../../core/services/community.service';
import { CreateCommunityRequest } from '../../../models/community.model';

@Component({
  selector: 'app-create-community-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './create-community-dialog.component.html',
  styleUrls: ['./create-community-dialog.component.scss']
})
export class CreateCommunityDialogComponent {
  communityForm: FormGroup;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private dialogRef: MatDialogRef<CreateCommunityDialogComponent>,
    private communityService: CommunityService,
    private snackBar: MatSnackBar
  ) {
    this.communityForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(21), Validators.pattern(/^[a-zA-Z0-9_]+$/)]],
      description: ['', [Validators.required, Validators.maxLength(500)]]
    });
  }

  onSubmit(): void {
    if (this.communityForm.invalid) {
      return;
    }

    this.isSubmitting = true;
    const communityData: CreateCommunityRequest = this.communityForm.value;

    this.communityService.createCommunity(communityData).subscribe({
      next: (community) => {
        this.snackBar.open(`Community c/${community.name} created!`, 'Close', { duration: 3000 });
        this.dialogRef.close(community);
      },
      error: (error) => {
        this.isSubmitting = false;
        const errorMessage = error.error?.error || 'Failed to create community';
        this.snackBar.open(errorMessage, 'Close', { duration: 4000 });
      }
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
