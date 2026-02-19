import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatBadgeModule } from '@angular/material/badge';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';

import { AuthService } from '../../core/services/auth.service';
import { User } from '../../models/user.model';
import { CreatePostDialogComponent } from '../../features/posts/create-post-dialog/create-post-dialog.component';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatMenuModule,
    MatBadgeModule,
    MatTooltipModule,
    MatDividerModule,
    MatDialogModule
  ],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  currentUser: User | null = null;
  searchQuery = '';
  notificationCount = 0;
  notifications: any[] = [];

  constructor(
    private authService: AuthService,
    private router: Router,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
    });

    // Mock notifications for demo
    this.notifications = [
      {
        icon: 'thumb_up',
        message: 'Someone upvoted your post "Best study spots on campus"',
        time: '2 hours ago'
      }
    ];
    this.notificationCount = this.notifications.length;
  }

  getJdenticonUrl(hash: string | undefined): string {
    if (!hash) return '';
    return `https://api.dicebear.com/7.x/identicon/svg?seed=${hash}`;
  }

  onSearch(): void {
    if (this.searchQuery.trim()) {
      this.router.navigate(['/search'], { 
        queryParams: { q: this.searchQuery } 
      });
    }
  }

  openCreatePost(): void {
    this.dialog.open(CreatePostDialogComponent, {
      width: '600px',
      maxWidth: '90vw'
    });
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
