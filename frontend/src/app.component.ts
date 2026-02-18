import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [],
  template: `
    <div style="padding: 50px; text-align: center; font-family: Arial;">
      <h1 style="color: #ff4500;">✅ Angular is Working!</h1>
      <p>If you see this, Angular is running correctly.</p>
      <p>Backend API: <a href="http://localhost:8080/health" target="_blank">Test Backend</a></p>
    </div>
  `
})
export class AppComponent {}