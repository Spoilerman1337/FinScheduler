export interface TagDto {
    id?: string;
    name?: string;
    isActive?: boolean;
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

export type TagStatusFilter = 'All' | 'Active' | 'Inactive';

export interface TagModification {
    name: string;
    isActive: boolean;
}
