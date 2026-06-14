import {afterEach, describe, expect, it, vi} from 'vitest';
import {API_BASE_URL} from '../config/api.ts';
import TagsService, {buildTagFilter} from './tags.ts';

describe('tags api', () => {
    afterEach(() => {
        vi.unstubAllGlobals();
    });

    it('buildTagFilter maps page, search, and active status to TagFilter', () => {
        // Arrange
        const params = {
            page: 4,
            pageSize: 25,
            searchTerm: 'groceries',
            statusFilter: 'Active' as const,
        };

        // Act
        const filter = buildTagFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 3,
            pageSize: 25,
            name: 'groceries',
            isActive: true,
        });
    });

    it('buildTagFilter omits optional criteria when they are inactive', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
        };

        // Act
        const filter = buildTagFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
        });
    });

    it('getDetailedInfo requests the tag by id endpoint and returns the tag', async () => {
        // Arrange
        const service = new TagsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(JSON.stringify({id: 'tag-1', name: 'Food'}), {
                status: 200,
                headers: {'Content-Type': 'application/json'},
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        const tag = await service.getDetailedInfo('tag-1');

        // Assert
        expect(fetchMock).toHaveBeenCalledWith(`${API_BASE_URL}/tags/tag-1`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        expect(tag).toEqual({id: 'tag-1', name: 'Food'});
    });

    it('getDetailedInfo returns null when the tag endpoint responds with 404', async () => {
        // Arrange
        const service = new TagsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(null, {
                status: 404,
                statusText: 'Not Found',
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        const tag = await service.getDetailedInfo('missing-tag');

        // Assert
        expect(tag).toBeNull();
    });
});
