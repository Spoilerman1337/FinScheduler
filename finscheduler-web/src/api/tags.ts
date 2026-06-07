import type {
    Lookup,
    PaginatedList,
} from './types';
import type {
    TagDto,
    TagFilter,
    TagLookupFilter,
    TagModification,
    TagStatusFilter,
} from './tags.types.ts';
import {FinschedulerApiClient} from './finscheduler-api-client.ts';

export function buildTagFilter(params: {
    page: number;
    pageSize: number;
    searchTerm: string;
    statusFilter: TagStatusFilter;
}): TagFilter {
    const filter: TagFilter = {
        page: params.page - 1,
        pageSize: params.pageSize,
    };

    if (params.searchTerm) {
        filter.name = params.searchTerm;
    }

    if (params.statusFilter !== 'All') {
        filter.isActive = params.statusFilter === 'Active';
    }

    return filter;
}

export default class TagsService extends FinschedulerApiClient {
    async getTags(filter?: TagFilter): Promise<PaginatedList<TagDto>> {
        const queryString = filter ? this.buildQueryString(filter) : '';
        const response = await fetch(`${this.baseUrl}/tags${queryString}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch tags: ${response.statusText}`);
        }

        return response.json();
    }

    async getTag(id: string): Promise<TagDto | null> {
        const response = await fetch(`${this.baseUrl}/tags/${id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.status === 404) {
            return null;
        }

        if (!response.ok) {
            throw new Error(`Failed to fetch tag: ${response.statusText}`);
        }

        return response.json();
    }

    async createTag(item: TagModification): Promise<string> {
        const response = await fetch(`${this.baseUrl}/tags`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: item.name,
                isActive: item.isActive,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to create tag: ${response.statusText}`);
        }

        return response.json();
    }

    async updateTag(id: string, item: TagModification): Promise<void> {
        const response = await fetch(`${this.baseUrl}/tags/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: item.name,
                isActive: item.isActive,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to update item: ${response.statusText}`);
        }
    }

    async getLookup(filter?: TagLookupFilter): Promise<PaginatedList<Lookup>> {
        const queryString = filter ? this.buildQueryString(filter) : '';
        const response = await fetch(`${this.baseUrl}/tags/lookup${queryString}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch tags: ${response.statusText}`);
        }

        return response.json();
    }
}
