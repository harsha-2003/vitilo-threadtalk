import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AuthService } from './auth.service';

describe('AuthService Unit Tests', () => {
  let service: AuthService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AuthService]
    });

    service = TestBed.inject(AuthService);
    localStorage.clear();
  });

  it('isAuthenticated() should return false when no token exists', () => {
  expect(service.isAuthenticated()).toBe(false);
});

it('isAuthenticated() should return true when token exists', () => {
  localStorage.setItem('token', 'fake-token');
  expect(service.isAuthenticated()).toBe(true);
});

  it('getCurrentUser() should return null when no user exists', () => {
    expect(service.getCurrentUser()).toBeNull();
  });

  it('getCurrentUser() should return user when user exists in localStorage', () => {
    const mockUser = { id: 1, email: 'a@b.com', anonymous_username: 'user1', avatar_hash: 'x' };
    localStorage.setItem('user', JSON.stringify(mockUser));
    expect(service.getCurrentUser()).toEqual(mockUser as any);
  });

  it('logout() should remove token and user from localStorage', () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({ id: 1 }));
    service.logout();
    expect(localStorage.getItem('token')).toBeNull();
    expect(localStorage.getItem('user')).toBeNull();
  });
});
