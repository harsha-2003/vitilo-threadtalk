export interface Community {
  id: number;
  name: string;
  description: string;
  icon_url: string;
  member_count: number;
  is_member: boolean;
  created_at: string;
}

export interface CreateCommunityRequest {
  name: string;
  description: string;
  icon_url?: string;
}
