import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { PostService } from '../../../core/services/post.service';
import { CommunityService } from '../../../core/services/community.service';
// import { HeaderComponent } from '../../../components/header/header.component';


interface Post {
  id: number;
  title: string;
  content: string;
  community_id: number;
  community_name: string;
  anonymous_username: string;
  vote_count: number;
  comment_count: number;
  created_at: string;
  user_vote: number;
}

interface Community {
  id: number;
  name: string;
  description: string;
  member_count: number;
}

@Component({
  selector: 'app-post-feed',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule], // Add HeaderComponent
  templateUrl: './post-feed.component.html',
  styleUrls: ['./post-feed.component.css']
})
export class PostFeedComponent implements OnInit {
  posts: Post[] = [];
  communities: Community[] = [];
  loading = false;
  showCreateModal = false;
  
  // Form data
  newPost = {
    title: '',
    content: '',
    community_id: 0,
    post_type: 'text' 
  };

  constructor(
    private postService: PostService,
    private communityService: CommunityService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadPosts();
    this.loadCommunities();
  }

  loadPosts(): void {
    this.loading = true;
    this.postService.getPosts().subscribe({
      next: (response: any) => {
        this.posts = response.posts || [];
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading posts:', error);
        this.loading = false;
      }
    });
  }

  loadCommunities(): void {
    this.communityService.getCommunities().subscribe({
      next: (communities) => {
        this.communities = communities;
      },
      error: (error) => {
        console.error('Error loading communities:', error);
      }
    });
  }

  openCreateModal(): void {
    this.showCreateModal = true;
  }

  closeCreateModal(): void {
    this.showCreateModal = false;
    this.newPost = { title: '', content: '', community_id: 0, post_type: 'text' };
  }

  createPost(): void {
    if (!this.newPost.title || !this.newPost.content || !this.newPost.community_id) {
      return;
    }

    this.postService.createPost(this.newPost).subscribe({
      next: (post) => {
        this.posts.unshift(post);
        this.closeCreateModal();
      },
      error: (error) => {
        console.error('Error creating post:', error);
      }
    });
  }

  votePost(postId: number, value: number): void {
    this.postService.votePost(postId, value).subscribe({
      next: (response: any) => {
        const post = this.posts.find(p => p.id === postId);
        if (post) {
          post.vote_count = response.vote_count;
          post.user_vote = value;
        }
      },
      error: (error) => {
        console.error('Error voting:', error);
      }
    });
  }

  formatVotes(votes: number): string {
    if (votes >= 1000) {
      return (votes / 1000).toFixed(1) + 'k';
    }
    return votes.toString();
  }

  getTimeAgo(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return 'just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + 'h ago';
    return Math.floor(seconds / 86400) + 'd ago';
  }

  goToPost(postId: number): void {
    this.router.navigate(['/post', postId]);
  }
}
