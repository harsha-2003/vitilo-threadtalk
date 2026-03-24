import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule, FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';

import { CommunityService } from '../../../core/services/community.service';

interface Community {
  id: number;
  name: string;
  description: string;
  member_count: number;
  is_member: boolean;
  icon_url?: string;
}

@Component({
  selector: 'app-community-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTabsModule
  ],
  templateUrl: './community-list.component.html',
  styleUrls: ['./community-list.component.scss']
})
export class CommunityListComponent implements OnInit {
  allCommunities: Community[] = [];
  myCommunities: Community[] = [];
  filteredCommunities: Community[] = [];
  searchQuery = '';
  isLoading = true;
  selectedTab = 0;
  showCreateDialog = false;
  isSubmitting = false;
  createForm: FormGroup;

  constructor(
    private communityService: CommunityService,
    private snackBar: MatSnackBar,
    private fb: FormBuilder
  ) {
    this.createForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(50)]],
      description: ['', [Validators.required, Validators.minLength(10), Validators.maxLength(500)]]
    });
  }

  ngOnInit(): void {
    this.loadCommunities();
    this.loadMyCommunities();
  }

  loadCommunities(): void {
    this.isLoading = true;
    this.communityService.getCommunities().subscribe({
      next: (response: any) => {
        this.allCommunities = (response.communities || response || []).sort((a: Community, b: Community) => 
          b.member_count - a.member_count
        );
        this.filteredCommunities = this.allCommunities;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Failed to load communities:', error);
        this.snackBar.open('Failed to load communities', 'Close', { duration: 3000 });
        this.isLoading = false;
      }
    });
  }

  loadMyCommunities(): void {
    this.communityService.getUserCommunities().subscribe({
      next: (response: any) => {
        this.myCommunities = response.communities || response || [];
      },
      error: (error) => {
        console.error('Failed to load your communities:', error);
      }
    });
  }

  onSearch(): void {
    const query = this.searchQuery.toLowerCase().trim();
    
    if (!query) {
      this.filteredCommunities = this.allCommunities;
      return;
    }

    this.filteredCommunities = this.allCommunities.filter(community => 
      community.name.toLowerCase().includes(query) ||
      community.description?.toLowerCase().includes(query)
    );
  }

  onTabChange(index: number): void {
    this.selectedTab = index;
    this.searchQuery = '';
    
    if (index === 0) {
      this.filteredCommunities = this.allCommunities;
    }
  }

  toggleJoin(community: Community): void {
    if (community.is_member) {
      this.leaveCommunity(community);
    } else {
      this.joinCommunity(community);
    }
  }

  joinCommunity(community: Community): void {
    this.communityService.joinCommunity(community.id).subscribe({
      next: () => {
        community.is_member = true;
        community.member_count++;
        this.snackBar.open(`Joined c/${community.name}`, 'Close', { duration: 2000 });
        this.loadMyCommunities();
      },
      error: (error) => {
        this.snackBar.open('Failed to join community', 'Close', { duration: 3000 });
      }
    });
  }

  leaveCommunity(community: Community): void {
    if (confirm(`Are you sure you want to leave c/${community.name}?`)) {
      this.communityService.leaveCommunity(community.id).subscribe({
        next: () => {
          community.is_member = false;
          community.member_count--;
          this.snackBar.open(`Left c/${community.name}`, 'Close', { duration: 2000 });
          this.loadMyCommunities();
        },
        error: (error) => {
          this.snackBar.open('Failed to leave community', 'Close', { duration: 3000 });
        }
      });
    }
  }

  openCreateCommunityDialog(): void {
    this.showCreateDialog = true;
  }

  closeCreateDialog(): void {
    this.showCreateDialog = false;
    this.createForm.reset();
  }

  createCommunity(): void {
    if (this.createForm.invalid) {
      return;
    }

    this.isSubmitting = true;

    const data = {
      name: this.createForm.value.name.trim(),
      description: this.createForm.value.description.trim()
    };

    this.communityService.createCommunity(data).subscribe({
      next: (community) => {
        this.allCommunities.unshift(community);
        this.filteredCommunities = this.allCommunities;
        this.closeCreateDialog();
        this.isSubmitting = false;
        this.snackBar.open('Community created! 🎉', 'Close', { duration: 3000 });
        this.loadCommunities();
        this.loadMyCommunities();
      },
      error: (error) => {
        console.error('Failed to create community:', error);
        this.isSubmitting = false;
        this.snackBar.open('Failed to create community', 'Close', { duration: 3000 });
      }
    });
  }

  getCommunityIcon(community: Community): string {
    if (community.icon_url) {
      return community.icon_url;
    }
    return `https://api.dicebear.com/7.x/initials/svg?seed=${community.name}`;
  }
}
