export interface Job {
    id: string;
    strategy_id: string;
    status: string;
}

export interface ListResponse {
    status: string;
    message: string;
    jobs: Job[];
}

export interface CreateRequest {
    strategy_id: string;
}

export interface CreateResponse {
    status: string;
    message: string;
    job: Job;
}