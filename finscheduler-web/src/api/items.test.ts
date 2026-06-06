import {afterEach, describe, expect, it, vi} from 'vitest';
import {API_BASE_URL} from '../config/api.ts';
import ItemsService, {buildItemFilter} from './items.ts';

describe('items api', () => {
    afterEach(() => {
        vi.unstubAllGlobals();
    });

    it('buildItemFilter maps page, filters, and valid numeric bounds to ItemFilter', () => {
        // Arrange
        const params = {
            page: 3,
            pageSize: 25,
            searchTerm: 'coffee',
            statusFilter: 'Active' as const,
            priceFrom: '100.50',
            priceTo: '300',
        };

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 2,
            pageSize: 25,
            name: 'coffee',
            isActive: true,
            priceFrom: 100.5,
            priceTo: 300,
        });
    });

    it('buildItemFilter skips inactive criteria when values are empty or invalid', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
            priceFrom: 'abc',
            priceTo: '',
        };

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
        });
    });

    it('getItem requests the item by id endpoint and returns the item', async () => {
        // Arrange
        const service = new ItemsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(JSON.stringify({id: 'item-1', name: 'Coffee'}), {
                status: 200,
                headers: {'Content-Type': 'application/json'},
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        const item = await service.getItem('item-1');

        // Assert
        expect(fetchMock).toHaveBeenCalledWith(`${API_BASE_URL}/items/item-1`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        expect(item).toEqual({id: 'item-1', name: 'Coffee'});
    });

    it('getItem returns null when the item endpoint responds with 404', async () => {
        // Arrange
        const service = new ItemsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(null, {
                status: 404,
                statusText: 'Not Found',
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        const item = await service.getItem('missing-item');

        // Assert
        expect(item).toBeNull();
    });
});
