// src/types/auth.types.ts
export interface LoginCredentials {
    username: string;
    password: string;
}

export interface TokenPayload {
    id: number;
    username: string;
    roles: string[];
    permissions: Permission[];
}

export interface Permission {
    name: string;
    modoule: string;
}