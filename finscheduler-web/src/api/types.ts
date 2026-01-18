export interface ItemDto {
    id?: string;
    name?: string;
    price?: number;
    description?: string;
    isActive?: boolean;
    createdAt?: string;
    updatedAt?: string | null;
    cashback?: number;
    category?: string;
}

export interface PaginatedList<T> {
    data: T[];
    count: number;
}

export interface ItemFilter {
    ids?: string[];
    name?: string;
    priceFrom?: number;
    priceTo?: number;
    description?: string;
    isActive?: boolean;
    createdFrom?: string;
    createdTo?: string;
    updatedFrom?: string;
    updatedTo?: string;
    cashbackFrom?: number;
    cashbackTo?: number;
    categories?: string[];
    page?: number;
    pageSize?: number;
}

