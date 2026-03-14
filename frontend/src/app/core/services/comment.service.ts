import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { Comment, CreateCommentRequest } from '../../models/comment.model';
import { environment } from '../../environments/environment';


@Injectable({
  providedIn: 'root'
})
export class CommentService {
   
  constructor(private http: HttpClient) {}

  getComments(postId: number): Observable<Comment[]> {console.log(`📤 Fetching comments for post ${postId}`);
    return this.http.get<Comment[]>(`${environment.apiUrl}/posts/${postId}/comments`).pipe(
      map(response => {
        console.log('📥 Comments response:', response);
        // Handle both array and object with comments property
        if (Array.isArray(response)) {
          return response;
        }
        // If backend returns {comments: [...]}
        return (response as any).comments || [];
      })
    );
  }

  createComment(request: CreateCommentRequest): Observable<Comment> {
    return this.http.post<Comment>(`${environment.apiUrl}/comments`, request);
  }

  deleteComment(id: number): Observable<any> {
    return this.http.delete(`${environment.apiUrl}/comments/${id}`);
  }

  voteComment(commentId: number, value: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/comments/${commentId}/vote`, { value });
  }
}
