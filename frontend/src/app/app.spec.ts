// frontend/src/app/app.spec.ts
import { TestBed } from '@angular/core/testing';
import { AppComponent } from '../app.component';

describe('AppComponent', () => {
  it('should create the app', async () => {
    await TestBed.configureTestingModule({
      imports: [AppComponent], // AppComponent is standalone in Angular 15+
    }).compileComponents();

    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
