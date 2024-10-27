// src/services/field/FieldService.ts
import apiClient from '@/utils/apiClient';
import type { ListResponse, CreateRequest, CreateResponse } from '@/types/FieldTypes';

class FieldService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/extractor/fields');
    return response.data;
  }

  static async create(payload: CreateRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>('/extractor/fields', payload);
    return response.data;
  }
}

export default FieldService;