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
            cashbackFrom: '5',
            cashbackTo: '15',
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
            cashbackFrom: 5,
            cashbackTo: 15,
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
            cashbackFrom: '',
            cashbackTo: 'oops',
        };

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
        });
    });

    it('buildItemFilter normalizes reversed price bounds before sending them to the API', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
            priceFrom: '300',
            priceTo: '100',
            cashbackFrom: '',
            cashbackTo: '',
        };

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
            priceFrom: 100,
            priceTo: 300,
        });
    });

    it('buildItemFilter normalizes reversed cashback bounds before sending them to the API', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
            priceFrom: '',
            priceTo: '',
            cashbackFrom: '25',
            cashbackTo: '10',
        };

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
            cashbackFrom: 10,
            cashbackTo: 25,
        });
    });

    it('buildItemFilter normalizes reversed created date bounds and expands them to full days', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
            dateFilter: {
                mode: 'created' as const,
                from: '2026-02-15',
                to: '2026-02-10',
            },
            priceFrom: '',
            priceTo: '',
            cashbackFrom: '',
            cashbackTo: '',
        };

        const expectedCreatedFrom = new Date('2026-02-10T00:00:00.000').toISOString();
        const expectedCreatedTo = new Date('2026-02-15T23:59:59.999').toISOString();

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
            createdFrom: expectedCreatedFrom,
            createdTo: expectedCreatedTo,
        });
    });

    it('buildItemFilter routes the selected updated date bounds to updated filter fields', () => {
        // Arrange
        const params = {
            page: 1,
            pageSize: 10,
            searchTerm: '',
            statusFilter: 'All' as const,
            dateFilter: {
                mode: 'updated' as const,
                from: '2026-03-02',
                to: '2026-03-05',
            },
            priceFrom: '',
            priceTo: '',
            cashbackFrom: '',
            cashbackTo: '',
        };

        const expectedUpdatedFrom = new Date('2026-03-02T00:00:00.000').toISOString();
        const expectedUpdatedTo = new Date('2026-03-05T23:59:59.999').toISOString();

        // Act
        const filter = buildItemFilter(params);

        // Assert
        expect(filter).toEqual({
            page: 0,
            pageSize: 10,
            updatedFrom: expectedUpdatedFrom,
            updatedTo: expectedUpdatedTo,
        });
    });

    it('getDetailedInfo requests the item by id endpoint and returns the item', async () => {
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
        const item = await service.getDetailedInfo('item-1');

        // Assert
        expect(fetchMock).toHaveBeenCalledWith(`${API_BASE_URL}/items/item-1`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        expect(item).toEqual({id: 'item-1', name: 'Coffee'});
    });

    it('getDetailedInfo returns null when the item endpoint responds with 404', async () => {
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
        const item = await service.getDetailedInfo('missing-item');

        // Assert
        expect(item).toBeNull();
    });

    it('updateCashbackByTag sends a PATCH request with the selected tag and cashback', async () => {
        // Arrange
        const service = new ItemsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(null, {
                status: 204,
                statusText: 'No Content',
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        await service.updateCashbackByTag('tag-1', 17);

        // Assert
        expect(fetchMock).toHaveBeenCalledWith(`${API_BASE_URL}/items/cashback/tag`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                tagId: 'tag-1',
                cashback: 17,
            }),
        });
    });

    it('updateCashbackByItems sends a PATCH request with selected item ids and cashback', async () => {
        // Arrange
        const service = new ItemsService();
        const fetchMock = vi.fn().mockResolvedValue(
            new Response(null, {
                status: 204,
                statusText: 'No Content',
            }),
        );

        vi.stubGlobal('fetch', fetchMock);

        // Act
        await service.updateCashbackByItems(['item-1', 'item-2'], 23);

        // Assert
        expect(fetchMock).toHaveBeenCalledWith(`${API_BASE_URL}/items/cashback/items`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                itemIds: ['item-1', 'item-2'],
                cashback: 23,
            }),
        });
    });
});
