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
    tags?: Lookup[];
}

export interface TagDto {
    id?: string;
    name?: string;
    isActive?: boolean;
}

export interface PaginatedList<T> {
    data: T[];
    count: number;
}

export interface Lookup {
    value?: string;
    label?: string;
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

export interface TagFilter {
    ids?: string[];
    name?: string;
    isActive?: boolean;
    page?: number;
    pageSize?: number;
}

export interface TagLookupFilter {
    name?: string;
    page?: number;
    pageSize?: number;
}

export interface ItemModification {
    name: string
    price: number
    description?: string
    isActive: boolean
    cashback: number
    category: string
    tagIds: string[]
}

export interface TagModification {
    name: string
    isActive: boolean
}
