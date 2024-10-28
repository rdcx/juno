// src/services/job/JobService.ts
import apiClient from '@/utils/apiClient';
import type { 
  ListResponse, 
  CreateRequest, 
  CreateResponse, 
} from '@/types/JobTypes';

class StrategyService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/extractor/jobs');
    return response.data;
  }

  static async create(payload: CreateRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>('/extractor/jobs', payload);
    return response.data;
  }
}

export default StrategyService;