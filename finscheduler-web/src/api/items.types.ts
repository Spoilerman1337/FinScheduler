import type {Lookup} from './types.ts';

export interface ItemDateFilterValue {
    mode: ItemDateFilterMode;
    from: string;
    to: string;
}

export interface ItemListingDto {
    id?: string;
    name?: string;
    price?: number;
    isActive?: boolean;
    updatedAt?: string | null;
    cashback?: number;
}

export interface ItemDetailedDto {
    name?: string;
    price?: number;
    description?: string;
    isActive?: boolean;
    cashback?: number;
    category?: string;
    tags?: Lookup[];
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
