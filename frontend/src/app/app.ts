import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';

import { NavbarComponent } from '../app/components/navbar/navbar.component';
import { SidebarComponent } from '../app/components/sidebar/sidebar.component';
import { RightSidebarComponent } from '../app/components/right-sidebar/right-sidebar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,       // *ngIf
    RouterOutlet,       // <router-outlet>
    NavbarComponent,    // <app-navbar>
    SidebarComponent,   // <app-sidebar>
    RightSidebarComponent // <app-right-sidebar>
  ],
  templateUrl: './app.html',
})
export class AppComponent {
  isAuthenticated = false; // toggle for testing
}
