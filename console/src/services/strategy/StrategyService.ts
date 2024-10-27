// src/services/strategy/StrategyService.ts
import apiClient from '@/utils/apiClient';
import type { 
  ListResponse, 
  CreateRequest, 
  CreateResponse, 
  AddSelectorRequest, 
  AddSelectorResponse, 
  RemoveSelectorRequest, 
  RemoveSelectorResponse, 
  AddFieldRequest, 
  AddFieldResponse,
  RemoveFieldRequest,
  RemoveFieldResponse, 
  AddFilterRequest, 
  AddFilterResponse,
  RemoveFilterRequest,
  RemoveFilterResponse 
} from '@/types/StrategyTypes';

class StrategyService {
  static async list(): Promise<ListResponse> {
    const response = await apiClient.get<ListResponse>('/extractor/strategies');
    return response.data;
  }

  static async create(payload: CreateRequest): Promise<CreateResponse> {
    const response = await apiClient.post<CreateResponse>('/extractor/strategies', payload);
    return response.data;
  }

  static async AddSelector(payload: AddSelectorRequest): Promise<AddSelectorResponse> {
    const response = await apiClient.post<AddSelectorResponse>('/extractor/strategies/' + payload.strategy_id + '/selectors', payload);
    return response.data;
  }

  static async RemoveSelector(payload: RemoveSelectorRequest): Promise<RemoveSelectorResponse> {
    const response = await apiClient.delete<RemoveSelectorResponse>('/extractor/strategies/' + payload.strategy_id + '/selectors', { data: payload });
    return response.data;
  }

  static async AddField(payload: AddFieldRequest): Promise<AddFieldResponse> {
    const response = await apiClient.post<AddFieldResponse>('/extractor/strategies/' + payload.strategy_id + '/fields', payload);
    return response.data;
  }

  static async RemoveField(payload: RemoveFieldRequest): Promise<RemoveFieldResponse> {
    const response = await apiClient.delete<RemoveFieldResponse>('/extractor/strategies/' + payload.strategy_id + '/fields', { data: payload });
    return response.data;
  }

  static async AddFilter(payload: AddFilterRequest): Promise<AddFilterResponse> {
    const response = await apiClient.post<AddFilterResponse>('/extractor/strategies/' + payload.strategy_id + '/filters', payload);
    return response.data;
  }

  static async RemoveFilter(payload: RemoveFilterRequest): Promise<RemoveFilterResponse> {
    const response = await apiClient.delete<RemoveFilterResponse>('/extractor/strategies/' + payload.strategy_id + '/filters', { data: payload });
    return response.data;
  }
}

export default StrategyService;