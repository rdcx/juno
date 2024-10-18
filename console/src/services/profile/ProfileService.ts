// src/api/auth/AuthService.ts
import apiClient from '@/utils/apiClient';
import type { ProfileResponse } from '@/types/ProfileTypes';

class ProfileService {
  static async getProfile(): Promise<ProfileResponse> {
    const response = await apiClient.get<ProfileResponse>('/profile');
    return response.data;
  }
}

export default ProfileService;