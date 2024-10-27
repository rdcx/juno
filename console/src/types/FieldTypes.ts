export interface Field {
    id: string;
    selector_id: string;
    name: string;
    type: string;
}

export interface ListResponse {
    status: string;
    message: string;
    fields: Field[];
}

export interface CreateRequest {
    name: string;
    selector_id: string;
    visibility: string;
    type: string;
}

export interface CreateResponse {
    status: string;
    message: string;
    selector: Field;
}