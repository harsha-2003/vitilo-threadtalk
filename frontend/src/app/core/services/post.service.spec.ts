import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { PostService } from './post.service';
import { environment } from '../../environments/environment';

describe('PostService Unit Tests', () => {
  let service: PostService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [PostService]
    });

    service = TestBed.inject(PostService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('getUserPosts(userId) should call GET /users/:id/posts', () => {
    service.getUserPosts(2).subscribe();

    const req = httpMock.expectOne(`${environment.apiUrl}/users/2/posts`);
    expect(req.request.method).toBe('GET');
    req.flush({ posts: [] });
  });
});
