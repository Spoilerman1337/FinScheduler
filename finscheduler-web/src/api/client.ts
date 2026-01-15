import type { ItemDto, PaginatedList, ItemFilter } from './types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api';

class ApiClient {
    private baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    private buildQueryString(params: Record<string, any>): string {
        const searchParams = new URLSearchParams();
        
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                if (Array.isArray(value)) {
                    value.forEach(v => searchParams.append(key, String(v)));
                } else {
                    searchParams.append(key, String(value));
                }
            }
        });

        const queryString = searchParams.toString();
        return queryString ? `?${queryString}` : '';
    }

    async getItems(filter?: ItemFilter): Promise<PaginatedList<ItemDto>> {
        const queryString = filter ? this.buildQueryString(filter) : '';
        const response = await fetch(`${this.baseUrl}/items${queryString}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch items: ${response.statusText}`);
        }

        return response.json();
    }

    async getItemById(id: string): Promise<ItemDto> {
        const response = await fetch(`${this.baseUrl}/items/${id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch item: ${response.statusText}`);
        }

        return response.json();
    }

    async createItem(item: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>): Promise<string> {
        const response = await fetch(`${this.baseUrl}/items`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: item.name,
                price: item.price,
                description: item.description,
                isActive: item.isActive,
                cashback: item.cashback,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to create item: ${response.statusText}`);
        }

        return response.json();
    }

    async updateItem(id: string, item: Omit<ItemDto, 'id' | 'createdAt' | 'updatedAt'>): Promise<void> {
        const response = await fetch(`${this.baseUrl}/items/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: item.name,
                price: item.price,
                description: item.description,
                isActive: item.isActive,
                cashback: item.cashback,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to update item: ${response.statusText}`);
        }
    }

    async deleteItem(id: string): Promise<void> {
        const response = await fetch(`${this.baseUrl}/items/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error(`Failed to delete item: ${response.statusText}`);
        }
    }
}

export const apiClient = new ApiClient(API_BASE_URL);

