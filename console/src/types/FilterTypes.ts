export interface Filter {
    id: string;
    field_id: string;
    name: string;
    type: string;
    value: string;
}

export interface ListResponse {
    status: string;
    message: string;
    filters: Filter[];
}

export interface CreateRequest {
    name: string;
    field_id: string;
    value: string;
    type: string;
}

export interface CreateResponse {
    status: string;
    message: string;
    filter: Filter;
}