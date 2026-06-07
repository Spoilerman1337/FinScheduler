export interface ItemDateFilterValue {
    mode: ItemDateFilterMode;
    from: string;
    to: string;
}

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
    tags?: import('./types.ts').Lookup[];
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

export type ItemStatusFilter = 'All' | 'Active' | 'Inactive';

export type ItemDateFilterMode = 'created' | 'updated';

export interface ItemModification {
    name: string;
    price: number;
    description?: string;
    isActive: boolean;
    cashback: number;
    category: string;
    tagIds: string[];
}
