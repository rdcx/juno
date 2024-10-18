export interface BalanceResponse {
    status: string;
    message: string;
    balance: number;
}

export interface DepositResponse {
    status: string;
    message: string;
}

export interface DepositRequest {
    amount: number;
}