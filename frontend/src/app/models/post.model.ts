export interface Post {
  id: number;
  title: string;
  content: string;
  image_url: string;
  post_type: 'text' | 'image' | 'link';
  vote_count: number;
  comment_count: number;
  created_at: string;
  user_id: number;
  anonymous_username: string;
  avatar_hash: string;
  community_id: number;
  community_name: string;
  user_vote: number;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  community_id: number;
  post_type: string;
  image_url?: string;
}

export interface PostsResponse {
  posts: Post[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
