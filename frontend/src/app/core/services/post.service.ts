import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Post, CreatePostRequest, PostsResponse } from '../../models/post.model';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class PostService {
  constructor(private http: HttpClient) {}

  getPosts(page: number = 1, limit: number = 20, sort: string = 'new', communityId?: number): Observable<PostsResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('limit', limit.toString())
      .set('sort', sort);

    if (communityId) {
      params = params.set('community_id', communityId.toString());
    }

    return this.http.get<PostsResponse>(`${environment.apiUrl}/posts`, { params });
  }

  getPost(id: number): Observable<Post> {
    return this.http.get<Post>(`${environment.apiUrl}/posts/${id}`);
  }

  createPost(request: CreatePostRequest): Observable<Post> {
    return this.http.post<Post>(`${environment.apiUrl}/posts`, request);
  }

  deletePost(id: number): Observable<any> {
    return this.http.delete(`${environment.apiUrl}/posts/${id}`);
  }

  uploadImage(file: File): Observable<{ image_url: string }> {
    const formData = new FormData();
    formData.append('image', file);
    return this.http.post<{ image_url: string }>(`${environment.apiUrl}/posts/upload`, formData);
  }

  votePost(postId: number, value: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/posts/${postId}/vote`, { value });
  }
}
