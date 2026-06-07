import type {ItemDto, PaginatedList, ItemFilter, ItemModification} from './types';
import {FinschedulerApiClient} from './finscheduler-api-client.ts';

export type ItemStatusFilter = 'All' | 'Active' | 'Inactive';

function parsePriceValue(value: string): number | null {
    if (!value) {
        return null;
    }

    const normalized = value.replace(',', '.').trim();

    if (!normalized) {
        return null;
    }

    const parsed = Number(normalized);

    if (!Number.isFinite(parsed) || parsed < 0) {
        return null;
    }

    return parsed;
}

export function buildItemFilter(params: {
    page: number;
    pageSize: number;
    searchTerm: string;
    statusFilter: ItemStatusFilter;
    priceFrom: string;
    priceTo: string;
}): ItemFilter {
    const filter: ItemFilter = {
        page: params.page - 1,
        pageSize: params.pageSize,
    };

    if (params.searchTerm) {
        filter.name = params.searchTerm;
    }

    if (params.statusFilter !== 'All') {
        filter.isActive = params.statusFilter === 'Active';
    }

    const priceFrom = parsePriceValue(params.priceFrom);
    const priceTo = parsePriceValue(params.priceTo);

    if (priceFrom !== null && priceTo !== null) {
        filter.priceFrom = Math.min(priceFrom, priceTo);
        filter.priceTo = Math.max(priceFrom, priceTo);
        return filter;
    }

    if (priceFrom !== null) {
        filter.priceFrom = priceFrom;
    }

    if (priceTo !== null) {
        filter.priceTo = priceTo;
    }

    return filter;
}

export default class ItemsService extends FinschedulerApiClient {
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

    async getItem(id: string): Promise<ItemDto | null> {
        const response = await fetch(`${this.baseUrl}/items/${id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.status === 404) {
            return null;
        }

        if (!response.ok) {
            throw new Error(`Failed to fetch item: ${response.statusText}`);
        }

        return response.json();
    }

    async createItem(item: ItemModification): Promise<string> {
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
                category: item.category,
                cashback: item.cashback,
                tagIds: item.tagIds,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to create item: ${response.statusText}`);
        }

        return response.json();
    }

    async updateItem(id: string, item: ItemModification): Promise<void> {
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
                category: item.category,
                tagIds: item.tagIds,
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
