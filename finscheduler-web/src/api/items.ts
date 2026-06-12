import {FinschedulerApiClient} from './finscheduler-api-client.ts';
import type {PaginatedList} from './types.ts';
import type {
    ItemDateFilterValue,
    ItemDetailedDto,
    ItemFilter,
    ItemListingDto,
    ItemModification,
    ItemStatusFilter,
} from './items.types.ts';

const itemDateFilterFields = {
    created: {
        from: 'createdFrom',
        to: 'createdTo',
    },
    updated: {
        from: 'updatedFrom',
        to: 'updatedTo',
    },
} as const;

function applyDateFilterRange(filter: ItemFilter, dateFilter?: ItemDateFilterValue): void {
    if (!dateFilter) {
        return;
    }

    const dateRange = FinschedulerApiClient.buildDateRange(dateFilter.from, dateFilter.to);
    const fields = itemDateFilterFields[dateFilter.mode];

    if (dateRange.from !== null) {
        filter[fields.from] = FinschedulerApiClient.toLocalDayBoundaryIso(dateRange.from, false);
    }

    if (dateRange.to !== null) {
        filter[fields.to] = FinschedulerApiClient.toLocalDayBoundaryIso(dateRange.to, true);
    }
}

export function buildItemFilter(params: {
    page: number;
    pageSize: number;
    searchTerm: string;
    statusFilter: ItemStatusFilter;
    dateFilter?: ItemDateFilterValue;
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

    applyDateFilterRange(filter, params.dateFilter);

    const priceRange = FinschedulerApiClient.buildNonNegativeRange(
        params.priceFrom,
        params.priceTo,
    );

    if (priceRange.from !== null) {
        filter.priceFrom = priceRange.from;
    }

    if (priceRange.to !== null) {
        filter.priceTo = priceRange.to;
    }

    const cashbackRange = FinschedulerApiClient.buildNonNegativeRange(
        params.cashbackFrom,
        params.cashbackTo,
    );

    if (cashbackRange.from !== null) {
        filter.cashbackFrom = cashbackRange.from;
    }

    if (cashbackRange.to !== null) {
        filter.cashbackTo = cashbackRange.to;
    }

    return filter;
}

export default class ItemsService extends FinschedulerApiClient {
    async getListingInfo(filter?: ItemFilter): Promise<PaginatedList<ItemListingDto>> {
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

    async getDetailedInfo(id: string): Promise<ItemDetailedDto | null> {
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
