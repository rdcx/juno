// src/api/auth/AuthService.ts
import apiClient from '@/utils/apiClient';
import type { BalanceResponse, DepositResponse } from '@/types/TokenTypes';
import type { AxiosResponse } from 'axios';

class TokenService {
  static async balance(): Promise<BalanceResponse> {
    const response = await apiClient.get<BalanceResponse>('/tokens/balance');
    return response.data;
  }

  static async deposit(amount: number): Promise<DepositResponse> {

    try {
    const response = await apiClient.post<DepositResponse>('/tokens/deposit', { amount });

    return response.data;
    } catch (error: any) {
      return error.response.data; 
    }
  }
}

export default TokenService;