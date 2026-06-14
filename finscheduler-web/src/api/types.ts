export interface PaginatedList<T> {
    data: T[];
    count: number;
}

export interface Lookup {
    value: string;
    label?: string;
}
