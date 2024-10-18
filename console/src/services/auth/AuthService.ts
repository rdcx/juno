// src/api/auth/AuthService.ts
import apiClient from '@/utils/apiClient';
import type { LoginPayload, RegisterPayload, AuthResponse } from '@/types/AuthTypes';

class AuthService {
  static async login(payload: LoginPayload): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/token', payload);
    return response.data;
  }

  static async register(payload: RegisterPayload): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/users', payload);
    return response.data;
  }
}

export default AuthService;