// src/types/common.types.ts
export interface PaginationParams {
    page?: number;
    limit?: number;
}

export interface SearchParams exxtends PaginationParams {
    search?: string;
    sortBy?: string;
    sortOrder?: 'ASC' | 'DESC';
}

