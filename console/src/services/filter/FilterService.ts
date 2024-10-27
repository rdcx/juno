// src/services/filter/FilterService.ts
import apiClient from '@/utils/apiClient';
import type { ListResponse, CreateRequest, CreateResponse } from '@/types/FilterTypes';

class FilterService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/extractor/filters');
    return response.data;
  }

  static async create(payload: CreateRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>('/extractor/filters', payload);
    return response.data;
  }
}

export default FilterService;