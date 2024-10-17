// src/api/auth/AuthTypes.ts
export interface LoginPayload {
    email: string;
    password: string;
}

export interface RegisterPayload {
    name: string;
    email: string;
    password: string;
}

export interface AuthResponse {
    status: string;
    message: string;
    token: string;
}