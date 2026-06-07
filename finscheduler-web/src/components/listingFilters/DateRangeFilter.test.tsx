import {fireEvent, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import DateRangeFilter from './DateRangeFilter.tsx';

const singleDateMode = [
    {
        value: 'created',
        label: 'Created',
    },
] as const;

const multipleDateModes = [
    {
        value: 'created',
        label: 'Created',
    },
    {
        value: 'updated',
        label: 'Updated',
    },
] as const;

describe('DateRangeFilter', () => {
    it('normalizes reversed date bounds before applying them', async () => {
        // Arrange
        const onChange = vi.fn();
        const user = userEvent.setup();

        renderWithProviders(
            <DateRangeFilter
                label="Date"
                modes={singleDateMode}
                value={{mode: 'created', from: '', to: ''}}
                onChange={onChange}
                fromInputLabel="Created from"
                toInputLabel="Created to"
            />,
        );

        // Act
        await user.click(screen.getByRole('button', {name: 'Created'}));
        fireEvent.change(await screen.findByLabelText('Created from'), {
            target: {value: '18.06.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Created from'));
        fireEvent.change(screen.getByLabelText('Created to'), {
            target: {value: '12.06.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Created to'));
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(onChange).toHaveBeenCalledWith({
            mode: 'created',
            from: '2025-06-12',
            to: '2025-06-18',
        });
    });

    it('does not render the mode switch when only one mode is available', async () => {
        // Arrange
        const user = userEvent.setup();

        renderWithProviders(
            <DateRangeFilter
                label="Date"
                modes={singleDateMode}
                value={{mode: 'created', from: '', to: ''}}
                onChange={() => undefined}
            />,
        );

        // Act
        await user.click(screen.getByRole('button', {name: 'Created'}));

        // Assert
        expect(screen.getAllByRole('button', {name: 'Created'})).toHaveLength(1);
    });

    it('applies the selected date mode together with the date range', async () => {
        // Arrange
        const onChange = vi.fn();
        const user = userEvent.setup();

        renderWithProviders(
            <DateRangeFilter
                label="Date"
                modes={multipleDateModes}
                value={{mode: 'created', from: '', to: ''}}
                onChange={onChange}
                fromInputLabel="Updated from"
                toInputLabel="Updated to"
            />,
        );

        // Act
        await user.click(screen.getByRole('button', {name: 'Created'}));
        await user.click(screen.getByRole('button', {name: 'Updated'}));
        fireEvent.change(await screen.findByLabelText('Updated from'), {
            target: {value: '10.07.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Updated from'));
        fireEvent.change(screen.getByLabelText('Updated to'), {
            target: {value: '11.07.2025'},
        });
        fireEvent.blur(screen.getByLabelText('Updated to'));
        await user.click(screen.getByRole('button', {name: 'Применить'}));

        // Assert
        expect(onChange).toHaveBeenCalledWith({
            mode: 'updated',
            from: '2025-07-10',
            to: '2025-07-11',
        });
    });
});
