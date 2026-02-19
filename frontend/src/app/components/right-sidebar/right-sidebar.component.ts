import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

import { CommunityService } from '../../core/services/community.service';
import { Community } from '../../models/community.model';

@Component({
  selector: 'app-right-sidebar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatButtonModule,
    MatIconModule
  ],
  templateUrl: './right-sidebar.component.html',
  styleUrls: ['./right-sidebar.component.scss']
})
export class RightSidebarComponent implements OnInit {
  trendingCommunities: Community[] = [];

  constructor(private communityService: CommunityService) {}

  ngOnInit(): void {
    this.loadTrendingCommunities();
  }

  loadTrendingCommunities(): void {
    this.communityService.getCommunities().subscribe({
      next: (communities) => {
        this.trendingCommunities = communities
          .sort((a, b) => b.member_count - a.member_count)
          .slice(0, 5);
      },
      error: (error) => {
        console.error('Failed to load communities:', error);
      }
    });
  }

  toggleJoin(community: Community): void {
    if (community.is_member) {
      this.communityService.leaveCommunity(community.id).subscribe({
        next: () => {
          community.is_member = false;
          community.member_count--;
        }
      });
    } else {
      this.communityService.joinCommunity(community.id).subscribe({
        next: () => {
          community.is_member = true;
          community.member_count++;
        }
      });
    }
  }
}
