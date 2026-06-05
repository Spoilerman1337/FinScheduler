import {describe, expect, it} from 'vitest';
import {buildItemFilter} from './items.ts';

describe('items api', () => {
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
});
