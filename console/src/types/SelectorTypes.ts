export interface Selector {
    id: string;
    name: string;
    value: string;
    visibility: string;
}

export interface ListResponse {
    status: string;
    message: string;
    selectors: Selector[];
}

export interface CreateRequest {
    name: string;
    value: string;
    visibility: string;
}

export interface CreateResponse {
    status: string;
    message: string;
    selector: Selector;
}