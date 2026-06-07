import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import NumberRangeFilter from './NumberRangeFilter.tsx';

describe('NumberRangeFilter', () => {
    it('allows negative values when no minimum is provided', async () => {
        // Arrange
        const onChange = vi.fn();
        const user = userEvent.setup();

        renderWithProviders(
            <NumberRangeFilter
                label="Balance"
                value={{from: '', to: ''}}
                onChange={onChange}
                fromInputLabel="Balance from"
                toInputLabel="Balance to"
                applyLabel="Apply"
                clearLabel="Clear"
            />,
        );

        // Act
        await user.click(screen.getByRole('button', {name: 'Balance'}));
        await user.type(await screen.findByLabelText('Balance from'), '-5');
        await user.type(screen.getByLabelText('Balance to'), '10');
        await user.click(screen.getByRole('button', {name: 'Apply'}));

        // Assert
        expect(onChange).toHaveBeenCalledWith({from: '-5', to: '10'});
    });

    it('clamps negative values to zero when minimum is set to zero', async () => {
        // Arrange
        const onChange = vi.fn();
        const user = userEvent.setup();

        renderWithProviders(
            <NumberRangeFilter
                label="Price"
                value={{from: '', to: ''}}
                onChange={onChange}
                min={0}
                fromInputLabel="Price from"
                toInputLabel="Price to"
                applyLabel="Apply"
                clearLabel="Clear"
            />,
        );

        // Act
        await user.click(screen.getByRole('button', {name: 'Price'}));
        await user.type(await screen.findByLabelText('Price from'), '-5');
        await user.type(screen.getByLabelText('Price to'), '10');
        await user.click(screen.getByRole('button', {name: 'Apply'}));

        // Assert
        expect(onChange).toHaveBeenCalledWith({from: '0', to: '10'});
    });
});
