export interface Transaction {
    id: string;
    amount: number;
    key: string;
    meta: string;
}

export interface ListResponse {
    status: string;
    message: string;
    transactions: Transaction[];
}