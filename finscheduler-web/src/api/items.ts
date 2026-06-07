import type {ItemDto, PaginatedList, ItemFilter, ItemModification} from './types';
import {FinschedulerApiClient} from './finscheduler-api-client.ts';

export type ItemStatusFilter = 'All' | 'Active' | 'Inactive';

function parseNonNegativeNumberValue(value: string): number | null {
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

function buildNonNegativeRange(fromValue: string, toValue: string): {
    from: number | null;
    to: number | null;
} {
    const from = parseNonNegativeNumberValue(fromValue);
    const to = parseNonNegativeNumberValue(toValue);

    if (from !== null && to !== null && from > to) {
        return {
            from: to,
            to: from,
        };
    }

    return {from, to};
}

export function buildItemFilter(params: {
    page: number;
    pageSize: number;
    searchTerm: string;
    statusFilter: ItemStatusFilter;
    priceFrom: string;
    priceTo: string;
    cashbackFrom: string;
    cashbackTo: string;
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

    const priceRange = buildNonNegativeRange(params.priceFrom, params.priceTo);

    if (priceRange.from !== null) {
        filter.priceFrom = priceRange.from;
    }

    if (priceRange.to !== null) {
        filter.priceTo = priceRange.to;
    }

    const cashbackRange = buildNonNegativeRange(params.cashbackFrom, params.cashbackTo);

    if (cashbackRange.from !== null) {
        filter.cashbackFrom = cashbackRange.from;
    }

    if (cashbackRange.to !== null) {
        filter.cashbackTo = cashbackRange.to;
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
