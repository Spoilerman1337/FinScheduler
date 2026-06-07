import {describe, expect, it} from 'vitest';
import {mapLookupsToSelectOptions} from './shared.ts';

describe('features shared', () => {
    it('maps lookups to select options and skips entries without value', () => {
        // Arrange
        const lookups = [{value: 'tag-1', label: 'Food'}, {value: 'tag-2'}, {label: 'Broken'}];

        // Act
        const options = mapLookupsToSelectOptions(lookups);

        // Assert
        expect(options).toEqual([
            {value: 'tag-1', label: 'Food'},
            {value: 'tag-2', label: 'tag-2'},
        ]);
    });

    it('returns an empty array when lookups are missing', () => {
        // Arrange
        const lookups = undefined;

        // Act
        const options = mapLookupsToSelectOptions(lookups);

        // Assert
        expect(options).toEqual([]);
    });
});
