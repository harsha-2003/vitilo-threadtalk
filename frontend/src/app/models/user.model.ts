export interface User {
  id: number;
  email: string;
  anonymous_username: string;
  avatar_hash: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}