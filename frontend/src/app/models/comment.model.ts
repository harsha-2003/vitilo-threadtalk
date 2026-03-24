export interface Comment {
  id: number;
  content: string;
  vote_count: number;
  created_at: string;
  user_id: number;
  anonymous_username: string;
  avatar_hash: string;
  parent_id: number | null;
  post_id: number;
  replies: Comment[];
  user_vote: number;
}

export interface CreateCommentRequest {
  content: string;
  post_id: number;
  parent_id?: number|null;
}
