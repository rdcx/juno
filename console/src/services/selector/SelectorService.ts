// src/services/selector/SelectorService.ts
import apiClient from '@/utils/apiClient';
import type { ListResponse, CreateRequest, CreateResponse } from '@/types/SelectorTypes';

class SelectorService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/extractor/selectors');
    return response.data;
  }

  static async create(payload: CreateRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>('/extractor/selectors', payload);
    return response.data;
  }
}

export default SelectorService;