// src/api/auth/AuthService.ts
import apiClient from '@/utils/apiClient';
import type { ListResponse } from '@/types/TransactionTypes';

class TransactionService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/transactions');
    return response.data;
  }
}

export default TransactionService;