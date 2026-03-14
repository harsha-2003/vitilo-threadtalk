import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { debounceTime, distinctUntilChanged, Subject, Subscription } from 'rxjs';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  searchQuery = '';
  searchResults: any[] = [];
  showSearchResults = false;
  isLoggedIn = false;
  currentUser: User | null = null;
  
  private searchSubject = new Subject<string>();
  private userSubscription?: Subscription;

  constructor(
    private readonly router: Router,
    private readonly authService: AuthService
  ) {}

  ngOnInit(): void {
    // Subscribe to current user observable
    this.userSubscription = this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
      this.isLoggedIn = !!user;
    });

    // Setup search with debounce
    this.searchSubject
      .pipe(
        debounceTime(300),
        distinctUntilChanged()
      )
      .subscribe(query => {
        this.performSearch(query);
      });
  }

  ngOnDestroy(): void {
    // Clean up subscription
    if (this.userSubscription) {
      this.userSubscription.unsubscribe();
    }
  }

  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    if (!target) return;
    
    const query = target.value;
    this.searchQuery = query;
    
    if (query.trim().length > 0) {
      this.searchSubject.next(query);
      this.showSearchResults = true;
    } else {
      this.showSearchResults = false;
      this.searchResults = [];
    }
  }

  performSearch(query: string): void {
    if (!query.trim()) {
      this.searchResults = [];
      return;
    }

    console.log('Searching for:', query);

    // TODO: Replace with actual API call
    // For now, using mock data
    this.searchResults = [
      {
        type: 'post',
        id: 1,
        title: `Post about ${query}`,
        community: 'General'
      },
      {
        type: 'community',
        id: 1,
        name: query,
        members: 100
      }
    ];
  }

  onSearchSubmit(event: Event): void {
    event.preventDefault();
    
    if (this.searchQuery.trim()) {
      this.router.navigate(['/search'], { 
        queryParams: { q: this.searchQuery } 
      });
      this.closeSearchResults();
    }
  }

  selectSearchResult(result: any): void {
    if (result.type === 'post') {
      this.router.navigate(['/post', result.id]);
    } else if (result.type === 'community') {
      this.router.navigate(['/communities', result.id]);
    }
    this.closeSearchResults();
  }

  closeSearchResults(): void {
    this.showSearchResults = false;
    this.searchQuery = '';
    this.searchResults = [];
  }

  goToNotifications(): void {
    if (!this.isLoggedIn) {
      this.router.navigate(['/login']);
      return;
    }
    this.router.navigate(['/notifications']);
  }

  goToMessages(): void {
    if (!this.isLoggedIn) {
      this.router.navigate(['/login']);
      return;
    }
    this.router.navigate(['/messages']);
  }

  goToProfile(): void {
    if (this.isLoggedIn) {
      this.router.navigate(['/profile']);
    } else {
      this.router.navigate(['/login']);
    }
  }

  goToHome(): void {
    if (this.isLoggedIn) {
      this.router.navigate(['/feed']);
    } else {
      this.router.navigate(['/']);
    }
  }

  logout(): void {
    if (confirm('Are you sure you want to logout?')) {
      this.authService.logout();
      this.router.navigate(['/']);
    }
  }

  getUserInitial(): string {
    if (this.currentUser?.anonymous_username) {
      return this.currentUser.anonymous_username.charAt(0).toUpperCase();
    }
    return '?';
  }
}
