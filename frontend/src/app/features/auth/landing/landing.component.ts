import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div style="text-align: center; padding: 50px;">
      <h1>Welcome to Vitilo ThreadTalk</h1>
      <p>Anonymous campus discussions for UF students</p>
      <a routerLink="/register" style="padding: 10px 20px; background: #0079d3; color: white; text-decoration: none; border-radius: 4px;">
        Get Started
      </a>
      <br><br>
      <a routerLink="/login" style="padding: 10px 20px; background: #ff4500; color: white; text-decoration: none; border-radius: 4px;">
        Login
      </a>
    </div>
  `
})
export class LandingComponent {}
