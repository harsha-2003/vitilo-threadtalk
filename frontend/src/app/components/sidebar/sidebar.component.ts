import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog } from '@angular/material/dialog';

import { CommunityService } from '../../core/services/community.service';
import { Community } from '../../models/community.model';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatIconModule,
    MatButtonModule,
    MatDividerModule,
    MatTooltipModule
  ],
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss']
})
export class SidebarComponent implements OnInit {
  userCommunities: Community[] = [];

  constructor(
    private communityService: CommunityService,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadUserCommunities();
  }

  loadUserCommunities(): void {
    this.communityService.getUserCommunities().subscribe({
      next: (communities) => {
        this.userCommunities = communities;
      },
      error: (error) => {
        console.error('Failed to load communities:', error);
      }
    });
  }

  createCommunity(): void {
    // Will implement in Sprint 2
    console.log('Create community clicked');
  }
}
