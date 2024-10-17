
export interface User {
    id: string;
    email: string;
    name: string;
}

export interface ProfileResponse {
    status: string;
    message: string;
    user: User;
}