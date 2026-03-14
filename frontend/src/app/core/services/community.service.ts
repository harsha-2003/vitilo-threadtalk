import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Community, CreateCommunityRequest } from '../../models/community.model';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class CommunityService {
  constructor(private http: HttpClient) {}

  getCommunities(): Observable<Community[]> {
    return this.http.get<Community[]>(`${environment.apiUrl}/communities`);
  }

  getCommunity(id: number): Observable<Community> {
    return this.http.get<Community>(`${environment.apiUrl}/communities/${id}`);
  }

  createCommunity(request: CreateCommunityRequest): Observable<Community> {
    return this.http.post<Community>(`${environment.apiUrl}/communities`, request);
  }

  joinCommunity(id: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/communities/${id}/join`, {});
  }

  leaveCommunity(id: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/communities/${id}/leave`, {});
  }

  getUserCommunities(): Observable<Community[]> {
    return this.http.get<Community[]>(`${environment.apiUrl}/communities/user/joined`);
  }
}
