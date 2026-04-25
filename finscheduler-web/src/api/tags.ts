import type {Lookup, PaginatedList, TagDto, TagFilter, TagLookupFilter} from './types';
import {FinschedulerApiClient} from "./finscheduler-api-client.ts";

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

    async createTag(item: Omit<TagDto, 'id'>): Promise<string> {
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

    async updateTag(id: string, item: Omit<TagDto, 'id'>): Promise<void> {
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
};

